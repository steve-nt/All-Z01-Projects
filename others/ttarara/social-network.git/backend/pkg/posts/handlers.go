package posts

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/utils"
)

type postResponse struct {
	PostID   int       `json:"post_id"`
	UserID   int       `json:"user_id"`
	Author   string    `json:"author"`
	Content  string    `json:"content"`
	Privacy  string    `json:"privacy"`
	CreatedAt time.Time `json:"created_at"`

	ImageURL string `json:"image_url,omitempty"`
}

type commentResponse struct {
	CommentID int       `json:"comment_id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	ImageURL string `json:"image_url,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, map[string]string{
		"error":   code,
		"message": message,
	})
}

func viewerIDFromRequest(r *http.Request) int {
	c, err := r.Cookie("session")
	if err != nil || c == nil {
		return 0
	}
	if !utils.IsValidSession(c.Value) {
		return 0
	}
	return utils.GetUserIDFromSession(c.Value)
}

func requirePostPrivacyValue(p string) (string, error) {
	p = strings.TrimSpace(strings.ToLower(p))
	switch p {
	case "public", "almost_private", "private":
		return p, nil
	default:
		return "", errors.New("invalid privacy")
	}
}

func canViewerSeePost(db *sql.DB, viewerID int, postID int) (bool, error) {
	var authorID int
	var privacy string
	err := db.QueryRow(`SELECT user_id, privacy FROM Posts WHERE post_id = ?`, postID).Scan(&authorID, &privacy)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
		return false, err
	}

	if viewerID != 0 && viewerID == authorID {
		return true, nil
	}

	switch privacy {
	case "public":
		return true, nil
	case "almost_private":
		if viewerID == 0 {
			return false, nil
		}
		var ok bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM Followers
				WHERE follower_id = ? AND following_id = ? AND status = 'accepted'
			)
		`, viewerID, authorID).Scan(&ok)
		return ok, err
	case "private":
		if viewerID == 0 {
			return false, nil
		}
		var ok bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM Post_Visibility
				WHERE post_id = ? AND user_id = ?
			)
		`, postID, viewerID).Scan(&ok)
		return ok, err
	default:
		// Unknown privacy value => hide
		return false, nil
	}
}

// PostsHandler
// GET /api/posts?user_id=&limit=&offset=
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	db := sqlite.GetDB()
	viewerID := viewerIDFromRequest(r)

	var (
		limit  = 20
		offset = 0
	)
	if s := r.URL.Query().Get("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if s := r.URL.Query().Get("offset"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n >= 0 {
			offset = n
		}
	}

	var filterUserID *int
	if s := r.URL.Query().Get("user_id"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			filterUserID = &n
		} else {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "invalid user_id")
			return
		}
	}

	// Privacy-aware feed:
	// - viewer=0: only public
	// - viewer>0: own posts OR public OR almost_private (accepted follower) OR private (Post_Visibility)
	var (
		args  []any
		query strings.Builder
	)
	query.WriteString(`
		SELECT
			p.post_id,
			p.user_id,
			COALESCE(NULLIF(u.nickname, ''), u.email) AS author,
			p.content,
			p.privacy,
			p.created_at,
			COALESCE((
				SELECT pi.image_url
				FROM Posts_Images pi
				WHERE pi.post_id = p.post_id
				ORDER BY pi.image_id ASC
				LIMIT 1
			), '') AS image_url
		FROM Posts p
		JOIN Users u ON u.user_id = p.user_id
		WHERE
	`)

	if viewerID == 0 {
		// Guest: only public posts (regardless of profile privacy)
		query.WriteString(` p.privacy = 'public' `)
	} else {
		// Authenticated: post-level privacy rules only
		query.WriteString(`
			(
				p.user_id = ?
				OR p.privacy = 'public'
				OR (
					p.privacy = 'almost_private'
					AND EXISTS (
						SELECT 1 FROM Followers f
						WHERE f.follower_id = ? AND f.following_id = p.user_id AND f.status = 'accepted'
					)
				)
				OR (
					p.privacy = 'private'
					AND EXISTS (
						SELECT 1 FROM Post_Visibility pv
						WHERE pv.post_id = p.post_id AND pv.user_id = ?
					)
				)
			)
		`)
		args = append(args, viewerID, viewerID, viewerID)
	}

	if filterUserID != nil {
		query.WriteString(` AND p.user_id = ? `)
		args = append(args, *filterUserID)
	}

	query.WriteString(` ORDER BY p.created_at DESC LIMIT ? OFFSET ?`)
	args = append(args, limit, offset)

	rows, err := db.Query(query.String(), args...)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch posts")
		return
	}
	defer rows.Close()

	posts := make([]postResponse, 0)
	for rows.Next() {
		var p postResponse
		var imageURL string
		if err := rows.Scan(&p.PostID, &p.UserID, &p.Author, &p.Content, &p.Privacy, &p.CreatedAt, &imageURL); err != nil {
			continue
		}
		if imageURL != "" {
			p.ImageURL = imageURL
		}
		posts = append(posts, p)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"viewer_id": viewerID,
		"posts":     posts,
		"limit":     limit,
		"offset":    offset,
	})
}

// PostViewHandler
// GET /api/posts/view?post_id=
func PostViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "post_id is required")
		return
	}

	db := sqlite.GetDB()
	viewerID := viewerIDFromRequest(r)

	ok, err := canViewerSeePost(db, viewerID, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "not_found", "post not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check permissions")
		return
	}
	if !ok {
		writeJSONError(w, http.StatusForbidden, "forbidden", "you cannot view this post")
		return
	}

	var p postResponse
	var imageURL string
	err = db.QueryRow(`
		SELECT
			p.post_id,
			p.user_id,
			COALESCE(NULLIF(u.nickname, ''), u.email) AS author,
			p.content,
			p.privacy,
			p.created_at,
			COALESCE((
				SELECT pi.image_url
				FROM Posts_Images pi
				WHERE pi.post_id = p.post_id
				ORDER BY pi.image_id ASC
				LIMIT 1
			), '') AS image_url
		FROM Posts p
		JOIN Users u ON u.user_id = p.user_id
		WHERE p.post_id = ?
	`, postID).Scan(&p.PostID, &p.UserID, &p.Author, &p.Content, &p.Privacy, &p.CreatedAt, &imageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "not_found", "post not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch post")
		return
	}
	if imageURL != "" {
		p.ImageURL = imageURL
	}

	writeJSON(w, http.StatusOK, p)
}

// CreatePostHandler
// POST /api/posts/create
//
// JSON:
// {
//   "content": "text",
//   "privacy": "public" | "almost_private" | "private",
//   "visible_to": [2,3],           // required only when privacy=private
//   "image_filename": "123_...png" // optional; returned by /api/upload-image
// }
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil || cookie == nil || !utils.IsValidSession(cookie.Value) {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}
	userID := utils.GetUserIDFromSession(cookie.Value)
	if userID == 0 {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	var payload struct {
		Content       string `json:"content"`
		Privacy       string `json:"privacy"`
		VisibleTo     []int  `json:"visible_to"`
		ImageFilename string `json:"image_filename"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body")
		return
	}

	content := strings.TrimSpace(payload.Content)
	hasImage := strings.TrimSpace(payload.ImageFilename) != ""
	if content == "" && !hasImage {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "post must have text or an image")
		return
	}
	privacy, err := requirePostPrivacyValue(payload.Privacy)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "privacy must be public|almost_private|private")
		return
	}

	db := sqlite.GetDB()
	tx, err := db.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to start transaction")
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO Posts (user_id, content, privacy) VALUES (?, ?, ?)`, userID, content, privacy)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to create post")
		return
	}
	postID64, _ := res.LastInsertId()
	postID := int(postID64)

	// Private visibility list: must be accepted followers of the creator.
	if privacy == "private" {
		if len(payload.VisibleTo) == 0 {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "visible_to is required for private posts")
			return
		}

		// De-dup and validate membership (accepted follower)
		seen := make(map[int]struct{}, len(payload.VisibleTo))
		inserted := 0
		for _, v := range payload.VisibleTo {
			if v <= 0 || v == userID {
				continue
			}
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}

			var isFollower bool
			if err := tx.QueryRow(`
				SELECT EXISTS (
					SELECT 1 FROM Followers
					WHERE follower_id = ? AND following_id = ? AND status = 'accepted'
				)
			`, v, userID).Scan(&isFollower); err != nil {
				writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to validate visibility")
				return
			}
			if !isFollower {
				writeJSONError(w, http.StatusBadRequest, "invalid_request", "visible_to must contain only accepted followers")
				return
			}

			if _, err := tx.Exec(`
				INSERT INTO Post_Visibility (post_id, user_id, visibility)
				VALUES (?, ?, 'allowed')
			`, postID, v); err != nil {
				writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to save visibility list")
				return
			}
			inserted++
		}
		if inserted == 0 {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "visible_to must contain at least one accepted follower (not yourself)")
			return
		}
	}

	// Optional image linking: accept filename returned by /api/upload-image.
	// We store metadata best-effort from the saved file path.
	if strings.TrimSpace(payload.ImageFilename) != "" {
		filename := strings.TrimSpace(payload.ImageFilename)
		imageURL := "/frontend/uploads/images/" + filename

		// Best-effort infer file type from filename extension.
		fileType := "JPEG"
		switch strings.ToLower(strings.TrimPrefix(filepathExt(filename), ".")) {
		case "png":
			fileType = "PNG"
		case "gif":
			fileType = "GIF"
		case "jpg", "jpeg":
			fileType = "JPEG"
		default:
			// Keep default (JPEG) to satisfy CHECK constraint; upload handler already restricts types.
		}

		// Use filename as original_name (upload handler doesn't persist original name).
		// file_size is stored as 0 if stat fails; doesn't violate NOT NULL for INTEGER in SQLite.
		var fileSize int64
		if fi, statErr := statUploadedImage(filename); statErr == nil {
			fileSize = fi
		}

		if _, err := tx.Exec(`
			INSERT INTO Posts_Images (post_id, user_id, filename, original_name, file_size, file_type, image_type, image_url)
			VALUES (?, ?, ?, ?, ?, ?, 'post', ?)
		`, postID, userID, filename, filename, fileSize, fileType, imageURL); err != nil {
			// If this fails, still keep the post (image is optional)
			// but return a warning to the client.
			if err := tx.Commit(); err != nil {
				writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to finalize post")
				return
			}
			writeJSON(w, http.StatusCreated, map[string]any{
				"post_id": postID,
				"warning": "post created but image metadata failed to save",
			})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to commit")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"post_id": postID,
	})
}

// CommentsHandler
// GET  /api/posts/comments?post_id=
// POST /api/posts/comments  (auth) JSON: { "post_id": 1, "content": "...", "image_url": "/frontend/uploads/images/x.gif" }
func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listCommentsHandler(w, r)
	case http.MethodPost:
		createCommentHandler(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

func listCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "post_id is required")
		return
	}

	db := sqlite.GetDB()
	viewerID := viewerIDFromRequest(r)
	ok, err := canViewerSeePost(db, viewerID, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "not_found", "post not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check permissions")
		return
	}
	if !ok {
		writeJSONError(w, http.StatusForbidden, "forbidden", "you cannot view this post")
		return
	}

	rows, err := db.Query(`
		SELECT
			c.id,
			c.post_id,
			c.user_id,
			COALESCE(NULLIF(u.nickname, ''), u.email) AS author,
			c.content,
			c.created_at,
			COALESCE((
				SELECT ci.image_path
				FROM Comments_Images ci
				WHERE ci.comment_id = c.id
				ORDER BY ci.id ASC
				LIMIT 1
			), '') AS image_url
		FROM Comments c
		JOIN Users u ON u.user_id = c.user_id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to fetch comments")
		return
	}
	defer rows.Close()

	comments := make([]commentResponse, 0)
	for rows.Next() {
		var c commentResponse
		var imageURL string
		if err := rows.Scan(&c.CommentID, &c.PostID, &c.UserID, &c.Author, &c.Content, &c.CreatedAt, &imageURL); err != nil {
			continue
		}
		if imageURL != "" {
			c.ImageURL = imageURL
		}
		comments = append(comments, c)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"post_id":  postID,
		"comments": comments,
	})
}

func createCommentHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie == nil || !utils.IsValidSession(cookie.Value) {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}
	userID := utils.GetUserIDFromSession(cookie.Value)
	if userID == 0 {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	var payload struct {
		PostID   int    `json:"post_id"`
		Content  string `json:"content"`
		ImageURL string `json:"image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body")
		return
	}
	if payload.PostID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "post_id is required")
		return
	}
	content := strings.TrimSpace(payload.Content)
	if content == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "content is required")
		return
	}

	db := sqlite.GetDB()
	ok, err := canViewerSeePost(db, userID, payload.PostID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "not_found", "post not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to check permissions")
		return
	}
	if !ok {
		writeJSONError(w, http.StatusForbidden, "forbidden", "you cannot comment on this post")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to start transaction")
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO Comments (post_id, user_id, content) VALUES (?, ?, ?)`, payload.PostID, userID, content)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to create comment")
		return
	}
	commentID64, _ := res.LastInsertId()
	commentID := int(commentID64)

	if strings.TrimSpace(payload.ImageURL) != "" {
		imageURL := strings.TrimSpace(payload.ImageURL)
		imgType := inferImageTypeFromURL(imageURL)
		if imgType == "" {
			writeJSONError(w, http.StatusBadRequest, "invalid_request", "image_url must be JPEG|PNG|GIF")
			return
		}
		if _, err := tx.Exec(`
			INSERT INTO Comments_Images (comment_id, image_path, image_type)
			VALUES (?, ?, ?)
		`, commentID, imageURL, imgType); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to save comment image")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "db_error", "failed to commit")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"comment_id": commentID,
	})
}

func filepathExt(name string) string {
	i := strings.LastIndexByte(name, '.')
	if i < 0 {
		return ""
	}
	return name[i:]
}

func inferImageTypeFromURL(u string) string {
	ext := strings.ToLower(filepathExt(u))
	switch ext {
	case ".jpg", ".jpeg":
		return "JPEG"
	case ".png":
		return "PNG"
	case ".gif":
		return "GIF"
	default:
		return ""
	}
}

func statUploadedImage(filename string) (int64, error) {
	// Files are saved under backend/frontend/uploads/images/ by the existing upload handler.
	path := "frontend/uploads/images/" + filename
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

