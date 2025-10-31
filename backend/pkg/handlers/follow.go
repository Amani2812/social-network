package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/auth"
	"social-network/pkg/models"
	"strconv"
	"time"
)

type FollowHandler struct {
	followRepo *models.FollowRepository
	userRepo   *models.UserRepository
}

func NewFollowHandler(followRepo *models.FollowRepository, userRepo *models.UserRepository) *FollowHandler {
	return &FollowHandler{followRepo: followRepo, userRepo: userRepo}
}

type FollowRequest struct {
	FollowingID int `json:"following_id"`
}

type FollowResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (h *FollowHandler) SendFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.FollowingID == userID {
		http.Error(w, "Cannot follow yourself", http.StatusBadRequest)
		return
	}

	// Check if user exists
	followingUser, err := h.userRepo.GetByID(req.FollowingID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if already following or request exists
	existing, err := h.followRepo.GetFollowRequest(userID, req.FollowingID)
	if err == nil && existing != nil {
		http.Error(w, "Follow request already exists", http.StatusConflict)
		return
	}

	status := "pending"
	if followingUser.IsPublic {
		status = "accepted"
	}

	follow := &models.Follow{
		FollowerID:  userID,
		FollowingID: req.FollowingID,
		Status:      status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.followRepo.Create(follow); err != nil {
		http.Error(w, "Failed to send follow request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FollowResponse{
		Success: true,
		Message: "Follow request sent successfully",
	})
}

func (h *FollowHandler) RespondToFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	followIDStr := r.URL.Query().Get("follow_id")
	action := r.URL.Query().Get("action")

	followID, err := strconv.Atoi(followIDStr)
	if err != nil {
		http.Error(w, "Invalid follow ID", http.StatusBadRequest)
		return
	}

	follow, err := h.followRepo.GetByID(followID)
	if err != nil {
		http.Error(w, "Follow request not found", http.StatusNotFound)
		return
	}

	if follow.FollowingID != userID {
		http.Error(w, "Unauthorized to respond to this request", http.StatusForbidden)
		return
	}

	var status string
	switch action {
	case "accept":
		status = "accepted"
	case "reject":
		status = "rejected"
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err := h.followRepo.UpdateStatus(followID, status); err != nil {
		http.Error(w, "Failed to update follow request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FollowResponse{
		Success: true,
		Message: "Follow request " + action + "ed successfully",
	})
}

func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	followingIDStr := r.URL.Query().Get("following_id")
	followingID, err := strconv.Atoi(followingIDStr)
	if err != nil {
		http.Error(w, "Invalid following ID", http.StatusBadRequest)
		return
	}

	follow, err := h.followRepo.GetFollowRequest(userID, followingID)
	if err != nil {
		http.Error(w, "Follow relationship not found", http.StatusNotFound)
		return
	}

	if err := h.followRepo.Delete(follow.ID); err != nil {
		http.Error(w, "Failed to unfollow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FollowResponse{
		Success: true,
		Message: "Unfollowed successfully",
	})
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	followers, err := h.followRepo.GetFollowers(userID)
	if err != nil {
		http.Error(w, "Failed to get followers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	following, err := h.followRepo.GetFollowing(userID)
	if err != nil {
		http.Error(w, "Failed to get following", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}

func (h *FollowHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	requests, err := h.followRepo.GetPendingRequests(userID)
	if err != nil {
		http.Error(w, "Failed to get pending requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}
