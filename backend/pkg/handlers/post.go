package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/auth"
	"social-network/pkg/models"
	"strconv"
	"time"
)

type PostHandler struct {
	postRepo    *models.PostRepository
	commentRepo *models.CommentRepository
	userRepo    *models.UserRepository
}

func NewPostHandler(postRepo *models.PostRepository, commentRepo *models.CommentRepository, userRepo *models.UserRepository) *PostHandler {
	return &PostHandler{postRepo: postRepo, commentRepo: commentRepo, userRepo: userRepo}
}

type CreatePostRequest struct {
	Content   string  `json:"content"`
	ImagePath *string `json:"image_path"`
	Privacy   string  `json:"privacy"`
}

type PostResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Post    *models.Post `json:"post,omitempty"`
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	if req.Privacy != "public" && req.Privacy != "almost_private" && req.Privacy != "private" {
		http.Error(w, "Invalid privacy setting", http.StatusBadRequest)
		return
	}

	post := &models.Post{
		UserID:    userID,
		Content:   req.Content,
		ImagePath: req.ImagePath,
		Privacy:   req.Privacy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.postRepo.Create(post); err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PostResponse{
		Success: true,
		Message: "Post created successfully",
		Post:    post,
	})
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postRepo.GetByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Check privacy
	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok && post.Privacy != "public" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if post.Privacy == "private" && post.UserID != userID {
		// Check if user is in the private list (not implemented yet, assume not)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if post.Privacy == "almost_private" && post.UserID != userID {
		// Check if user follows the post owner
		followRepo := models.NewFollowRepository(h.userRepo.DB)
		isFollowing, err := followRepo.IsFollowing(userID, post.UserID)
		if err != nil || !isFollowing {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	targetUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	currentUserID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		currentUserID = 0 // Not logged in
	}

	posts, err := h.postRepo.GetByUserID(targetUserID, currentUserID)
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	posts, err := h.postRepo.GetFeed(userID)
	if err != nil {
		http.Error(w, "Failed to get feed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postRepo.GetByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if post.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	post.Content = req.Content
	post.ImagePath = req.ImagePath
	post.Privacy = req.Privacy
	post.UpdatedAt = time.Now()

	if err := h.postRepo.Update(post); err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PostResponse{
		Success: true,
		Message: "Post updated successfully",
		Post:    post,
	})
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postRepo.GetByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if post.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if err := h.postRepo.Delete(postID); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PostResponse{
		Success: true,
		Message: "Post deleted successfully",
	})
}
