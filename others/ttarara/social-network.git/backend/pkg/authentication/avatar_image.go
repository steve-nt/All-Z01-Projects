package authentication

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/utils"
	"strings"
	"time"

	"golang.org/x/image/draw"
)

const (
	MaxFileSize   = 20 * 1024 * 1024 // 20MB
	UploadDir     = "frontend/uploads/images"
	ThumbnailDir  = "frontend/uploads/thumbnails"
	ThumbnailSize = 300
)

// AvatarUploadHandler handles avatar image upload for user profiles
// Updates the avatar_path in the Users table
func AvatarUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check authentication
	cookie, err := r.Cookie("session")
	if err != nil || cookie == nil || !utils.IsValidSession(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := utils.GetUserIDFromSession(cookie.Value)
	if userID == 0 {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// Parse multipart form (32MB max memory)
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the file from form
	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileType, err := validateImageType(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file.Seek(0, 0)

	// Process avatar upload
	avatarPath, err := processAvatarUpload(file, fileHeader, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update avatar_path in Users table
	db := sqlite.GetDB()

	// Get old avatar path to delete it later
	var oldAvatarPath string
	db.QueryRow("SELECT avatar_path FROM Users WHERE user_id = ?", userID).Scan(&oldAvatarPath)

	// Update user's avatar_path
	_, err = db.Exec("UPDATE Users SET avatar_path = ? WHERE user_id = ?", avatarPath, userID)
	if err != nil {
		// Clean up file on database error
		filename := strings.TrimPrefix(avatarPath, "/frontend/uploads/images/")
		os.Remove(filepath.Join(UploadDir, filename))
		http.Error(w, "Error saving avatar", http.StatusInternalServerError)
		return
	}

	// Delete old avatar file if it exists
	if oldAvatarPath != "" && oldAvatarPath != avatarPath {
		oldPath := strings.TrimPrefix(oldAvatarPath, "/frontend/uploads/images/")
		os.Remove(filepath.Join(UploadDir, oldPath))
		// Also try to remove thumbnail
		os.Remove(filepath.Join(ThumbnailDir, oldPath))
	}

	// Return success response with avatar URL
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"success": true,
		"avatarPath": "%s",
		"fileSize": %d,
		"fileType": "%s"
	}`, avatarPath, fileHeader.Size, fileType)
}

// ImageUploadHandler handles image upload for posts
// NOTE: This is a placeholder for Part 3 (Posts & Groups)
// Currently saves files but doesn't store in Posts_Images table
// Part 3 will implement full post image storage functionality
func ImageUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check authentication
	cookie, err := r.Cookie("session")
	if err != nil || !utils.IsValidSession(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := utils.GetUserIDFromSession(cookie.Value)
	if userID == 0 {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// Parse multipart form (32MB max memory)
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the file from form
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size
	if fileHeader.Size > MaxFileSize {
		http.Error(w, "File too large. Maximum size is 20MB", http.StatusBadRequest)
		return
	}

	// Validate file type
	fileType, err := validateImageType(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Reset file pointer to beginning
	file.Seek(0, 0)

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%d_%s", userID, timestamp, sanitizeFilename(fileHeader.Filename))

	// Ensure upload directories exist
	os.MkdirAll(UploadDir, 0755)
	os.MkdirAll(ThumbnailDir, 0755)

	// Save original image
	imagePath := filepath.Join(UploadDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		os.Remove(imagePath) // Clean up on error
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Create thumbnail
	thumbnailPath, err := createThumbnail(imagePath, filename, fileType)
	if err != nil {
		// Log error but don't fail the upload
		fmt.Printf("Warning: Could not create thumbnail: %v\n", err)
	}

	// Note: This handler is for post images, which should be stored in Posts_Images table
	// This is handled by Part 3 (Posts & Groups), so this is a placeholder
	// For now, return the image URL so it can be used in posts
	imageURL := "/frontend/uploads/images/" + filename
	thumbnailURL := ""
	if thumbnailPath != "" {
		thumbnailURL = "/frontend/uploads/thumbnails/" + filename
	}

	// Return success response with image URL
	// The actual database insertion should be done when creating the post (Part 3)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"success": true,
		"filename": "%s",
		"imageUrl": "%s",
		"thumbnailUrl": "%s",
		"fileSize": %d,
		"fileType": "%s"
	}`, filename, imageURL, thumbnailURL, fileHeader.Size, fileType)
}

// processAvatarUpload handles the avatar file upload and returns the avatar path
// This function can be used both during registration and profile updates
// Returns: avatarPath (relative URL), error
func processAvatarUpload(file multipart.File, fileHeader *multipart.FileHeader, userID int) (string, error) {
	// Validate file size
	if fileHeader.Size > MaxFileSize {
		return "", fmt.Errorf("file too large. Maximum size is 20MB")
	}

	// Validate file type
	fileType, err := validateImageType(file)
	if err != nil {
		return "", err
	}

	// Reset file pointer to beginning
	file.Seek(0, 0)

	// Generate unique filename for avatar
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("avatar_%d_%d_%s", userID, timestamp, sanitizeFilename(fileHeader.Filename))

	// Ensure upload directories exist
	os.MkdirAll(UploadDir, 0755)
	os.MkdirAll(ThumbnailDir, 0755)

	// Save original image
	imagePath := filepath.Join(UploadDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		os.Remove(imagePath) // Clean up on error
		return "", fmt.Errorf("error saving file: %w", err)
	}

	// Create thumbnail (optional for avatars, but useful)
	_, err = createThumbnail(imagePath, filename, fileType)
	if err != nil {
		// Log error but don't fail the upload
		fmt.Printf("Warning: Could not create thumbnail: %v\n", err)
	}

	// Return the relative path
	return "/frontend/uploads/images/" + filename, nil
}

// validateImageType checks if the uploaded file is a valid image type
func validateImageType(file io.ReadSeeker) (string, error) {
	// Read first 512 bytes for type detection
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("error reading file")
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Detect content type
	contentType := http.DetectContentType(buffer)

	switch contentType {
	case "image/jpeg":
		return "JPEG", nil
	case "image/png":
		return "PNG", nil
	case "image/gif":
		return "GIF", nil
	default:
		return "", fmt.Errorf("unsupported file type. Only JPEG, PNG, and GIF are allowed")
	}
}

// createThumbnail creates a thumbnail version of the uploaded image
func createThumbnail(originalPath, filename, fileType string) (string, error) {
	// Open original image
	originalFile, err := os.Open(originalPath)
	if err != nil {
		return "", err
	}
	defer originalFile.Close()

	// Decode image based on type
	var img image.Image
	switch fileType {
	case "JPEG":
		img, err = jpeg.Decode(originalFile)
	case "PNG":
		img, err = png.Decode(originalFile)
	case "GIF":
		img, err = gif.Decode(originalFile)
	default:
		return "", fmt.Errorf("unsupported image type for thumbnail")
	}

	if err != nil {
		return "", err
	}

	// Calculate thumbnail dimensions (maintain aspect ratio)
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	var newWidth, newHeight int
	if width > height {
		newWidth = ThumbnailSize
		newHeight = (height * ThumbnailSize) / width
	} else {
		newHeight = ThumbnailSize
		newWidth = (width * ThumbnailSize) / height
	}

	// Create thumbnail image
	thumbnail := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.NearestNeighbor.Scale(thumbnail, thumbnail.Rect, img, bounds, draw.Over, nil)

	// Save thumbnail
	thumbnailPath := filepath.Join(ThumbnailDir, filename)
	thumbnailFile, err := os.Create(thumbnailPath)
	if err != nil {
		return "", err
	}
	defer thumbnailFile.Close()

	// Encode thumbnail (always save as JPEG for consistency)
	err = jpeg.Encode(thumbnailFile, thumbnail, &jpeg.Options{Quality: 80})
	if err != nil {
		os.Remove(thumbnailPath)
		return "", err
	}

	return thumbnailPath, nil
}

// sanitizeFilename removes potentially dangerous characters from filename
func sanitizeFilename(filename string) string {
	// Remove path separators and other dangerous characters
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, "..", "_")
	filename = strings.ReplaceAll(filename, " ", "_")

	// Limit filename length
	if len(filename) > 100 {
		ext := filepath.Ext(filename)
		name := filename[:100-len(ext)]
		filename = name + ext
	}

	return filename
}

// DeleteAvatarHandler handles avatar deletion
func DeleteAvatarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check authentication
	cookie, err := r.Cookie("session")
	if err != nil || cookie == nil || !utils.IsValidSession(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := utils.GetUserIDFromSession(cookie.Value)
	if userID == 0 {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	db := sqlite.GetDB()

	// Get current avatar path
	var avatarPath string
	err = db.QueryRow("SELECT avatar_path FROM Users WHERE user_id = ?", userID).Scan(&avatarPath)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Delete avatar file from filesystem if it exists
	if avatarPath != "" {
		filename := strings.TrimPrefix(avatarPath, "/frontend/uploads/images/")
		os.Remove(filepath.Join(UploadDir, filename))
		// Also try to remove thumbnail
		os.Remove(filepath.Join(ThumbnailDir, filename))
	}

	// Clear avatar_path in database
	_, err = db.Exec("UPDATE Users SET avatar_path = NULL WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Error deleting avatar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"success": true}`)
}
