package handlers

import (
	"errors"
	"forum-image-upload/internal/backend/models"
	"forum-image-upload/internal/backend/services"
	"forum-image-upload/internal/utils"
	"log"
	"net/http"
)

type PostHandler struct {
	PostService *services.PostService
}

func (h *PostHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp, err := h.PostService.List(r.Context(), r)
	if err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, resp, 200)
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp, err := h.PostService.FindPostbyID(r.Context(), r)
	if err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp.Post.ID == "" {
		utils.JsonResponse(w, errors.New("No Post Found"), http.StatusNotFound)
		return
	}

	utils.JsonResponse(w, resp, 200)
}

func (h *PostHandler) CreateComment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value("user").(models.User)
	if user.UUID == "" {
		utils.JsonResponse(w, "User Action from guest", http.StatusBadRequest)
		return

	}

	if err := h.PostService.CreateComment(r.Context(), r, user); err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, "Comment Successfully created", 200)
}
func (h *PostHandler) StorePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value("user").(models.User)
	if user.UUID == "" {
		utils.JsonResponse(w, error.Error(errors.New("User Action from guest")), http.StatusBadRequest)
		return

	}
	if err := h.PostService.CreatePost(r.Context(), r, user); err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, "Post created", 200)
}

func (h *PostHandler) EditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value("user").(models.User)
	if user.UUID == "" {
		utils.JsonResponse(w, error.Error(errors.New("User Action from guest")), http.StatusInternalServerError)
		return

	}
	if err := h.PostService.EditPost(r.Context(), r, user); err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, "Post created", 200)
}

func (h *PostHandler) EditComment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value("user").(models.User)
	if user.UUID == "" {
		utils.JsonResponse(w, "User Action from guest", http.StatusBadRequest)
		return

	}

	if err := h.PostService.EditComment(r.Context(), r, user); err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, "Comment Successfully created", 200)
}

func (h *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := r.Context().Value("user").(models.User)
	if user.UUID == "" {
		utils.JsonResponse(w, error.Error(errors.New("User Action from guest")), http.StatusInternalServerError)
		return

	}
	if err := h.PostService.LikePost(r.Context(), r, user); err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, "User Liked Post", 200)
}

func (h *PostHandler) DislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value("user").(models.User)
	if user.UUID == "" {
		utils.JsonResponse(w, error.Error(errors.New("User Action from guest")), http.StatusInternalServerError)
		return

	}

	if err := h.PostService.DislikePost(r.Context(), r, user); err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, "User Disliked Post", 200)
}

func (h *PostHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := r.Context().Value("user").(models.User)
	if user.UUID == "" {
		utils.JsonResponse(w, error.Error(errors.New("User Action from guest")), http.StatusInternalServerError)
		return

	}
	if err := h.PostService.LikeComment(r.Context(), r, user); err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, "User Liked Comment", 200)
}

func (h *PostHandler) DislikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value("user").(models.User)
	if user.UUID == "" {
		utils.JsonResponse(w, error.Error(errors.New("User Action from guest")), http.StatusInternalServerError)
		return

	}

	if err := h.PostService.DislikeComment(r.Context(), r, user); err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, "User Disliked Comment", 200)
}

func (h *PostHandler) RemovePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value("user").(models.User)
	if user.UUID == "" {
		utils.JsonResponse(w, error.Error(errors.New("User Action from guest")), http.StatusInternalServerError)
		return

	}

	if err := h.PostService.RemovePost(r.Context(), r, user); err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, "User Deleted Post", 200)
}

func (h *PostHandler) RemoveComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value("user").(models.User)
	if user.UUID == "" {
		utils.JsonResponse(w, error.Error(errors.New("User Action from guest")), http.StatusBadRequest)
		return

	}

	if err := h.PostService.RemoveComment(r.Context(), r, user); err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, "User Deleted Comment", 200)
}
