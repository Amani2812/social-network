package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"social-network/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo *models.Repository
	hub  interface {
		BroadcastNewPost(postID int, userID int)
		SendFollowStatusUpdate(userID int)
		SendNotificationToUser(userID int, notifType, content string, senderID int)
	}
}

func NewHandler(repo *models.Repository, hub interface{ BroadcastNewPost(postID int, userID int); SendFollowStatusUpdate(userID int); SendNotificationToUser(userID int, notifType, content string, senderID int) }) *Handler {
	return &Handler{
		repo: repo,
		hub:  hub,
	}
}

// Helper functions
func (h *Handler) getUserFromSession(r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	session, err := h.repo.GetSession(cookie.Value)
	if err != nil {
		return nil, err
	}

	return h.repo.GetUserByID(session.UserID)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// Auth handlers
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		DateOfBirth string `json:"date_of_birth"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.repo.CreateUser(req.Email, req.Password, req.FirstName, req.LastName, req.DateOfBirth)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	session, err := h.repo.CreateSession(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	respondJSON(w, http.StatusOK, user)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		h.repo.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	respondJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Convert is_private to is_public for frontend
	response := map[string]interface{}{
		"id":            user.ID,
		"email":         user.Email,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"date_of_birth": user.DateOfBirth,
		"avatar_path":   user.AvatarPath,
		"nickname":      user.Nickname,
		"about_me":      user.AboutMe,
		"is_public":     !user.IsPrivate,
		"is_private":    user.IsPrivate,
		"created_at":    user.CreatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

// User handlers
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Convert is_private to is_public for frontend
	response := map[string]interface{}{
		"id":            user.ID,
		"email":         user.Email,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"date_of_birth": user.DateOfBirth,
		"avatar_path":   user.AvatarPath,
		"nickname":      user.Nickname,
		"about_me":      user.AboutMe,
		"is_public":     !user.IsPrivate,
		"is_private":    user.IsPrivate,
		"created_at":    user.CreatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		respondError(w, http.StatusBadRequest, "Search query required")
		return
	}

	users, err := h.repo.SearchUsers(query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to search users")
		return
	}

	respondJSON(w, http.StatusOK, users)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		FirstName  *string `json:"first_name"`
		LastName   *string `json:"last_name"`
		AvatarPath *string `json:"avatar_path"`
		Nickname   *string `json:"nickname"`
		AboutMe    *string `json:"about_me"`
		IsPublic   *bool   `json:"is_public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Convert is_public to is_private for database
	var isPrivate bool
	if req.IsPublic != nil {
		isPrivate = !(*req.IsPublic)
	} else {
		// Keep current value if not provided
		isPrivate = user.IsPrivate
	}

	if err := h.repo.UpdateUserProfile(user.ID, req.AvatarPath, req.Nickname, req.AboutMe, isPrivate); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}

// Follow handlers
func (h *Handler) SendFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		FollowingID int `json:"following_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.CreateFollowRequest(user.ID, req.FollowingID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to send follow request")
		return
	}

	// Create notification for all follow requests (both public and private profiles)
	targetUser, _ := h.repo.GetUserByID(req.FollowingID)
	if targetUser != nil {
		var content string
		if targetUser.IsPrivate {
			content = fmt.Sprintf("%s %s wants to follow you", user.FirstName, user.LastName)
		} else {
			content = fmt.Sprintf("%s %s started following you", user.FirstName, user.LastName)
		}
		h.repo.CreateNotification(req.FollowingID, "follow_request", content, &user.ID)
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Follow request sent"})
}

func (h *Handler) RespondToFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		FollowerID int    `json:"follower_id"`
		Accept     bool   `json:"accept"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	status := "rejected"
	if req.Accept {
		status = "accepted"
	}

	if err := h.repo.UpdateFollowStatus(req.FollowerID, user.ID, status); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to respond to follow request")
		return
	}

	// Send real-time WebSocket notification to the follower
	h.hub.SendFollowStatusUpdate(req.FollowerID)

	respondJSON(w, http.StatusOK, map[string]string{"message": "Follow request updated"})
}

func (h *Handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Try to get following_id from query parameter first (for compatibility)
	followingIDStr := r.URL.Query().Get("following_id")
	var followingID int
	
	if followingIDStr != "" {
		// From query parameter
		var err error
		followingID, err = strconv.Atoi(followingIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid following ID")
			return
		}
	} else {
		// From request body
		var req struct {
			FollowingID int `json:"following_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		followingID = req.FollowingID
	}

	if err := h.repo.DeleteFollow(user.ID, followingID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to unfollow")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Unfollowed successfully"})
}

func (h *Handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	followers, err := h.repo.GetFollowers(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get followers")
		return
	}

	respondJSON(w, http.StatusOK, followers)
}

func (h *Handler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	following, err := h.repo.GetFollowing(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get following")
		return
	}

	respondJSON(w, http.StatusOK, following)
}

func (h *Handler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	requests, err := h.repo.GetPendingFollowRequests(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get pending requests")
		return
	}

	respondJSON(w, http.StatusOK, requests)
}

func (h *Handler) GetFollowStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	followerIDStr := r.URL.Query().Get("follower_id")
	followingIDStr := r.URL.Query().Get("following_id")

	followerID, err := strconv.Atoi(followerIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid follower ID")
		return
	}

	followingID, err := strconv.Atoi(followingIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid following ID")
		return
	}

	// Check if following
	isFollowing := false
	hasPendingRequest := false

	follow, err := h.repo.GetFollowRelationship(followerID, followingID)
	if err == nil && follow != nil {
		if follow.Status == "accepted" {
			isFollowing = true
		} else if follow.Status == "pending" {
			hasPendingRequest = true
		}
	}

	respondJSON(w, http.StatusOK, map[string]bool{
		"is_following":         isFollowing,
		"has_pending_request": hasPendingRequest,
	})
}

// Post handlers
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		Content   string  `json:"content"`
		Privacy   string  `json:"privacy"`
		ImagePath *string `json:"image_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	post, err := h.repo.CreatePost(user.ID, req.Content, req.Privacy, req.ImagePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	// Broadcast new post to all connected clients
	h.hub.BroadcastNewPost(post.ID, user.ID)

	respondJSON(w, http.StatusCreated, post)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	post, err := h.repo.GetPost(postID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Post not found")
		return
	}

	respondJSON(w, http.StatusOK, post)
}

func (h *Handler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	posts, err := h.repo.GetUserPosts(userID, user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get posts")
		return
	}

	respondJSON(w, http.StatusOK, posts)
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	posts, err := h.repo.GetFeed(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get feed")
		return
	}

	respondJSON(w, http.StatusOK, posts)
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		PostID  int    `json:"post_id"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.UpdatePost(req.PostID, user.ID, req.Content); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update post")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Post updated successfully"})
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	if err := h.repo.DeletePost(postID, user.ID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Post deleted successfully"})
}

// Comment handlers
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		PostID    int     `json:"post_id"`
		Content   string  `json:"content"`
		ImagePath *string `json:"image_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	comment, err := h.repo.CreateComment(req.PostID, user.ID, req.Content, req.ImagePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	respondJSON(w, http.StatusCreated, comment)
}

func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	comments, err := h.repo.GetPostComments(postID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get comments")
		return
	}

	respondJSON(w, http.StatusOK, comments)
}

// Upload handler
func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "File too large")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to get file")
		return
	}
	defer file.Close()

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		respondError(w, http.StatusBadRequest, "Invalid file type. Only JPEG, PNG, and GIF are allowed")
		return
	}

	// Create uploads directory
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%d%s", user.ID, time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Create file
	dst, err := os.Create(filepath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"path": "/uploads/" + filename,
	})
}

// Group handlers
func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	group, err := h.repo.CreateGroup(user.ID, req.Title, req.Description)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create group")
		return
	}

	respondJSON(w, http.StatusCreated, group)
}

func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	groupIDStr := r.URL.Query().Get("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	group, err := h.repo.GetGroup(groupID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Group not found")
		return
	}

	respondJSON(w, http.StatusOK, group)
}

func (h *Handler) GetUserGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	groups, err := h.repo.GetUserGroups(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get groups")
		return
	}

	respondJSON(w, http.StatusOK, groups)
}

func (h *Handler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	_, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	groups, err := h.repo.GetAllGroups()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get groups")
		return
	}

	respondJSON(w, http.StatusOK, groups)
}

func (h *Handler) InviteToGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		GroupID int `json:"group_id"`
		UserID  int `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.InviteToGroup(req.GroupID, req.UserID, user.ID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create notification
	group, _ := h.repo.GetGroup(req.GroupID)
	if group != nil {
		content := fmt.Sprintf("You've been invited to join %s", group.Title)
		h.repo.CreateNotification(req.UserID, "group_invite", content, &req.GroupID)
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Invitation sent"})
}

func (h *Handler) RespondToGroupInvite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		GroupID int  `json:"group_id"`
		Accept  bool `json:"accept"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.RespondToGroupInvite(req.GroupID, user.ID, req.Accept); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to respond to invite")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Response recorded"})
}

func (h *Handler) RequestToJoinGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		GroupID int `json:"group_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.RequestToJoinGroup(req.GroupID, user.ID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create notification for group admins
	group, _ := h.repo.GetGroup(req.GroupID)
	if group != nil {
		members, _ := h.repo.GetGroupMembers(req.GroupID)
		for _, member := range members {
			if member.Role == "admin" {
				content := fmt.Sprintf("%s %s wants to join %s", user.FirstName, user.LastName, group.Title)
				h.repo.CreateNotification(member.UserID, "group_join_request", content, &req.GroupID)
				// Send real-time notification via WebSocket
				h.hub.SendNotificationToUser(member.UserID, "group_join_request", content, user.ID)
			}
		}
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Join request sent"})
}

func (h *Handler) GetGroupJoinRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	groupIDStr := r.URL.Query().Get("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Check if user is admin of the group
	members, err := h.repo.GetGroupMembers(groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get group members")
		return
	}

	isAdmin := false
	for _, member := range members {
		if member.UserID == user.ID && member.Role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		respondError(w, http.StatusForbidden, "Only admins can view join requests")
		return
	}

	requests, err := h.repo.GetGroupJoinRequests(groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get join requests")
		return
	}

	respondJSON(w, http.StatusOK, requests)
}

func (h *Handler) RespondToJoinRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		GroupID int  `json:"group_id"`
		UserID  int  `json:"user_id"`
		Accept  bool `json:"accept"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.RespondToJoinRequest(req.GroupID, req.UserID, user.ID, req.Accept); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create notification for the requester
	group, _ := h.repo.GetGroup(req.GroupID)
	if group != nil {
		var content string
		if req.Accept {
			content = fmt.Sprintf("Your request to join %s has been accepted", group.Title)
		} else {
			content = fmt.Sprintf("Your request to join %s has been declined", group.Title)
		}
		h.repo.CreateNotification(req.UserID, "group_join_request", content, &req.GroupID)
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Response recorded"})
}

func (h *Handler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	groupIDStr := r.URL.Query().Get("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Check if user is a member of the group
	members, err := h.repo.GetGroupMembers(groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get group members")
		return
	}

	isMember := false
	for _, member := range members {
		if member.UserID == user.ID {
			isMember = true
			break
		}
	}

	if !isMember {
		respondError(w, http.StatusForbidden, "Only group members can view member list")
		return
	}

	// Get members with user details
	rows, err := h.repo.GetGroupMembersWithDetails(groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get group members")
		return
	}

	respondJSON(w, http.StatusOK, rows)
}

func (h *Handler) CreateGroupPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		GroupID   int     `json:"group_id"`
		Content   string  `json:"content"`
		ImagePath *string `json:"image_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	post, err := h.repo.CreateGroupPost(req.GroupID, user.ID, req.Content, req.ImagePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, post)
}

func (h *Handler) GetGroupPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	groupIDStr := r.URL.Query().Get("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	posts, err := h.repo.GetGroupPosts(groupID, user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, posts)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		GroupID     int     `json:"group_id"`
		Title       string  `json:"title"`
		Description *string `json:"description"`
		EventTime   string  `json:"event_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	eventTime, err := time.Parse(time.RFC3339, req.EventTime)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid event time format")
		return
	}

	event, err := h.repo.CreateEvent(req.GroupID, user.ID, req.Title, req.Description, eventTime)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create event")
		return
	}

	// Notify all group members
	members, _ := h.repo.GetGroupMembers(req.GroupID)
	group, _ := h.repo.GetGroup(req.GroupID)
	if group != nil {
		for _, member := range members {
			if member.UserID != user.ID {
				content := fmt.Sprintf("New event in %s: %s", group.Title, req.Title)
				h.repo.CreateNotification(member.UserID, "event_invite", content, &event.ID)
				// Send real-time notification via WebSocket
				h.hub.SendNotificationToUser(member.UserID, "event_invite", content, user.ID)
			}
		}
	}

	respondJSON(w, http.StatusCreated, event)
}

func (h *Handler) GetGroupEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	groupIDStr := r.URL.Query().Get("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	events, err := h.repo.GetGroupEventsWithUserResponse(groupID, user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get events")
		return
	}

	respondJSON(w, http.StatusOK, events)
}

func (h *Handler) RespondToEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		EventID  int    `json:"event_id"`
		Response string `json:"response"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.RespondToEvent(req.EventID, user.ID, req.Response); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to respond to event")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Response recorded"})
}

// Message handlers
func (h *Handler) GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	otherUserIDStr := r.URL.Query().Get("user_id")
	otherUserID, err := strconv.Atoi(otherUserIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	messages, err := h.repo.GetPrivateMessages(user.ID, otherUserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get messages")
		return
	}

	respondJSON(w, http.StatusOK, messages)
}

func (h *Handler) GetGroupMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	groupIDStr := r.URL.Query().Get("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	messages, err := h.repo.GetGroupMessages(groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get messages")
		return
	}

	respondJSON(w, http.StatusOK, messages)
}

func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	conversations, err := h.repo.GetConversations(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get conversations")
		return
	}

	respondJSON(w, http.StatusOK, conversations)
}

// Notification handlers
func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notifications, err := h.repo.GetUserNotifications(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get notifications")
		return
	}

	respondJSON(w, http.StatusOK, notifications)
}

func (h *Handler) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req struct {
		NotificationID int `json:"notification_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.MarkNotificationAsRead(req.NotificationID, user.ID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

func (h *Handler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	count, err := h.repo.GetUnreadNotificationCount(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get unread count")
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"count": count})
}
