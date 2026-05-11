package services

import (
	"context"
	"encoding/json"
	"forum-authentication/internal/backend/models"
	"forum-authentication/internal/backend/repository"
	"forum-authentication/internal/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

type PostService struct {
	posts repository.PostRepository
}

func NewPostService(pr repository.PostRepository) *PostService { return &PostService{posts: pr} }

func (s *PostService) List(ctx context.Context, r *http.Request) (resp models.PostResponse, err error) {

	var username string
	var categories []string
	user := ctx.Value("user").(models.User)
	if user.UUID == "" {
		username = "guest"
	} else {
		username = user.Username
	}

	postCategories := utils.FindParamsFromURL(r, "category")
	if postCategories == nil {
		categories = append(resp.Category, "all")
	} else {
		categories = make([]string, len(postCategories))
		copy(categories, postCategories)
	}

	allPosts, err := s.posts.ListByCategory(ctx, postCategories)
	if err != nil {
		return models.PostResponse{}, err
	}

	resp = models.PostResponse{
		Posts:    allPosts,
		Username: username,
		Category: categories,
	}
	return resp, nil
}

func (s *PostService) FindPostbyID(ctx context.Context, r *http.Request) (models.PostByIdResponse, error) {
	var username string
	var response models.PostByIdResponse

	user := ctx.Value("user").(models.User)
	if user.UUID == "" {
		username = "guest"
	} else {
		username = user.Username
	}

	post_id := strings.TrimPrefix(r.URL.Path, "/postbyid/")
	post, err := s.posts.FindPostbyID(ctx, post_id)
	if err != nil {
		return response, err
	}
	response.Username = username
	response.Post = post

	return response, nil
}

func (s *PostService) CreateComment(ctx context.Context, r *http.Request, user models.User) error {
	post_id := strings.TrimPrefix(r.URL.Path, "/create-comment/")

	decoder := json.NewDecoder(r.Body)
	var commentContent string
	if err := decoder.Decode(&commentContent); err != nil {
		return err
	}

	return s.posts.CreateComment(ctx, commentContent, user.UUID, post_id, time.Now())
}
func (s *PostService) CreatePost(ctx context.Context, r *http.Request, user models.User) error {
	decoder := json.NewDecoder(r.Body)
	var postInfo models.PostInfo
	if err := decoder.Decode(&postInfo); err != nil {
		return err
	}

	return s.posts.CreatePost(ctx, &postInfo, user.UUID, time.Now().Unix())
}

func (s *PostService) EditPost(ctx context.Context, r *http.Request, user models.User) error {
	post_id := strings.TrimPrefix(r.URL.Path, "/edit-post/")
	decoder := json.NewDecoder(r.Body)
	var postInfo models.PostInfo
	if err := decoder.Decode(&postInfo); err != nil {
		return err
	}

	return s.posts.EditPost(ctx, &postInfo, post_id, user.UUID, time.Now().Unix())
}

func (s *PostService) EditComment(ctx context.Context, r *http.Request, user models.User) error {
	comment_id := strings.TrimPrefix(r.URL.Path, "/edit-comment/")
	decoder := json.NewDecoder(r.Body)
	var content string
	if err := decoder.Decode(&content); err != nil {
		return err
	}

	return s.posts.EditComment(ctx, content, comment_id, user.UUID, time.Now().Unix())
}

func (s *PostService) LikePost(ctx context.Context, r *http.Request, user models.User) error {
	post_id := strings.TrimPrefix(r.URL.Path, "/like-post/")
	err := s.posts.LikeButton(ctx, post_id, user, "post")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *PostService) DislikePost(ctx context.Context, r *http.Request, user models.User) error {
	post_id := strings.TrimPrefix(r.URL.Path, "/dislike-post/")
	err := s.posts.DislikeButton(ctx, post_id, user, "post")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *PostService) LikeComment(ctx context.Context, r *http.Request, user models.User) error {
	post_id := strings.TrimPrefix(r.URL.Path, "/like-comment/")
	err := s.posts.LikeButton(ctx, post_id, user, "comment")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *PostService) DislikeComment(ctx context.Context, r *http.Request, user models.User) error {
	post_id := strings.TrimPrefix(r.URL.Path, "/dislike-comment/")
	err := s.posts.DislikeButton(ctx, post_id, user, "comment")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *PostService) RemovePost(ctx context.Context, r *http.Request, user models.User) error {
	post_id := strings.TrimPrefix(r.URL.Path, "/remove-post/")
	err := s.posts.RemovePost(ctx, user.UUID, post_id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *PostService) RemoveComment(ctx context.Context, r *http.Request, user models.User) error {
	comment_id := strings.TrimPrefix(r.URL.Path, "/remove-comment/")
	err := s.posts.RemoveComment(ctx, user.UUID, comment_id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
