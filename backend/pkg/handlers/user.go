package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/auth"
	"social-network/pkg/models"
	"strconv"
	"strings"
)

type UserHandler struct {
	userRepo *models.UserRepository
}

func NewUserHandler(userRepo *models.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	userIDStr := parts[len(parts)-1]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if current user can view this profile
	currentUserID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		// Not logged in, only show public profiles
		if !user.IsPublic {
			http.Error(w, "Profile is private", http.StatusForbidden)
			return
		}
	} else {
		// Logged in, check if they can view private profiles
		if !user.IsPublic && currentUserID != userID {
			// Check if current user follows this user
		followRepo := models.NewFollowRepository(h.userRepo.DB)
			isFollowing, err := followRepo.IsFollowing(currentUserID, userID)
			if err != nil || !isFollowing {
				http.Error(w, "Profile is private", http.StatusForbidden)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updateData struct {
		Nickname    *string `json:"nickname"`
		AboutMe     *string `json:"about_me"`
		AvatarPath  *string `json:"avatar_path"`
		IsPublic    *bool   `json:"is_public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get current user
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update fields if provided
	if updateData.Nickname != nil {
		user.Nickname = updateData.Nickname
	}
	if updateData.AboutMe != nil {
		user.AboutMe = updateData.AboutMe
	}
	if updateData.AvatarPath != nil {
		user.AvatarPath = updateData.AvatarPath
	}
	if updateData.IsPublic != nil {
		user.IsPublic = *updateData.IsPublic
	}

	if err := h.userRepo.Update(user); err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter required", http.StatusBadRequest)
		return
	}

	users, err := h.userRepo.SearchByName(query)
	if err != nil {
		http.Error(w, "Failed to search users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
