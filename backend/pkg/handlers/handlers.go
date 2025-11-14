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
}

func NewHandler(repo *models.Repository) *Handler {
	return &Handler{repo: repo}
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

	respondJSON(w, http.StatusOK, user)
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

	respondJSON(w, http.StatusOK, user)
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
		AvatarPath *string `json:"avatar_path"`
		Nickname   *string `json:"nickname"`
		AboutMe    *string `json:"about_me"`
		IsPrivate  bool    `json:"is_private"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.UpdateUserProfile(user.ID, req.AvatarPath, req.Nickname, req.AboutMe, req.IsPrivate); err != nil {
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

	// Create notification
	targetUser, _ := h.repo.GetUserByID(req.FollowingID)
	if targetUser != nil && targetUser.IsPrivate {
		content := fmt.Sprintf("%s %s wants to follow you", user.FirstName, user.LastName)
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

	var req struct {
		FollowingID int `json:"following_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.DeleteFollow(user.ID, req.FollowingID); err != nil {
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

	groupIDStr := r.URL.Query().Get("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	events, err := h.repo.GetGroupEvents(groupID)
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
