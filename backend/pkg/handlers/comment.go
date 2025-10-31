package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/auth"
	"social-network/pkg/models"
	"strconv"
	"time"
)

type CommentHandler struct {
	commentRepo *models.CommentRepository
	postRepo    *models.PostRepository
}

func NewCommentHandler(commentRepo *models.CommentRepository, postRepo *models.PostRepository) *CommentHandler {
	return &CommentHandler{commentRepo: commentRepo, postRepo: postRepo}
}

type CreateCommentRequest struct {
	PostID    int     `json:"post_id"`
	Content   string  `json:"content"`
	ImagePath *string `json:"image_path"`
}

type CommentResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Comment *models.Comment  `json:"comment,omitempty"`
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	// Check if post exists and user can comment
	post, err := h.postRepo.GetByID(req.PostID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Check if user can view the post (same logic as in GetPost)
	if post.Privacy == "private" && post.UserID != userID {
		http.Error(w, "Cannot comment on private post", http.StatusForbidden)
		return
	}

	if post.Privacy == "almost_private" && post.UserID != userID {
		followRepo := models.NewFollowRepository(h.postRepo.DB)
		isFollowing, err := followRepo.IsFollowing(userID, post.UserID)
		if err != nil || !isFollowing {
			http.Error(w, "Cannot comment on this post", http.StatusForbidden)
			return
		}
	}

	comment := &models.Comment{
		PostID:    req.PostID,
		UserID:    userID,
		Content:   req.Content,
		ImagePath: req.ImagePath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.commentRepo.Create(comment); err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CommentResponse{
		Success: true,
		Message: "Comment created successfully",
		Comment: comment,
	})
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if user can view the post
	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		userID = 0
	}

	post, err := h.postRepo.GetByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if post.Privacy == "private" && post.UserID != userID {
		http.Error(w, "Cannot view comments", http.StatusForbidden)
		return
	}

	if post.Privacy == "almost_private" && post.UserID != userID {
		followRepo := models.NewFollowRepository(h.postRepo.DB)
		isFollowing, err := followRepo.IsFollowing(userID, post.UserID)
		if err != nil || !isFollowing {
			http.Error(w, "Cannot view comments", http.StatusForbidden)
			return
		}
	}

	comments, err := h.commentRepo.GetByPostID(postID)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	commentIDStr := r.URL.Query().Get("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	comment, err := h.commentRepo.GetByID(commentID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	if comment.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	comment.Content = req.Content
	comment.ImagePath = req.ImagePath
	comment.UpdatedAt = time.Now()

	if err := h.commentRepo.Update(comment); err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CommentResponse{
		Success: true,
		Message: "Comment updated successfully",
		Comment: comment,
	})
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	commentIDStr := r.URL.Query().Get("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	comment, err := h.commentRepo.GetByID(commentID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	if comment.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if err := h.commentRepo.Delete(commentID); err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CommentResponse{
		Success: true,
		Message: "Comment deleted successfully",
	})
}
