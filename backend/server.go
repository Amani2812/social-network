package main

import (
	"log"
	"net/http"

	"social-network/pkg/db"
	"social-network/pkg/handlers"
	"social-network/pkg/models"
	"social-network/pkg/websocket"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Initialize database
	database, err := db.NewDB("./social_network.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repository
	repo := models.NewRepository(database.DB)

	// Initialize WebSocket hub
	hub := websocket.NewHub(repo)
	go hub.Run()

	// Initialize handlers with hub
	handler := handlers.NewHandler(repo, hub)

	// Serve static files (uploads)
	fs := http.FileServer(http.Dir("./uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	// Auth routes
	http.HandleFunc("/api/auth/register", enableCORS(handler.Register))
	http.HandleFunc("/api/auth/login", enableCORS(handler.Login))
	http.HandleFunc("/api/auth/logout", enableCORS(handler.Logout))
	http.HandleFunc("/api/auth/me", enableCORS(handler.GetCurrentUser))

	// User routes
	http.HandleFunc("/api/users/", enableCORS(handler.GetUserByID))
	http.HandleFunc("/api/users/search", enableCORS(handler.SearchUsers))
	http.HandleFunc("/api/users/profile", enableCORS(handler.UpdateProfile))

	// Follow routes
	http.HandleFunc("/api/follow/request", enableCORS(handler.SendFollowRequest))
	http.HandleFunc("/api/follow/respond", enableCORS(handler.RespondToFollowRequest))
	http.HandleFunc("/api/follow/unfollow", enableCORS(handler.Unfollow))
	http.HandleFunc("/api/follow/status", enableCORS(handler.GetFollowStatus))
	http.HandleFunc("/api/follow/followers", enableCORS(handler.GetFollowers))
	http.HandleFunc("/api/follow/following", enableCORS(handler.GetFollowing))
	http.HandleFunc("/api/follow/pending", enableCORS(handler.GetPendingRequests))

	// Post routes
	http.HandleFunc("/api/posts", enableCORS(handler.CreatePost))
	http.HandleFunc("/api/posts/get", enableCORS(handler.GetPost))
	http.HandleFunc("/api/posts/user", enableCORS(handler.GetUserPosts))
	http.HandleFunc("/api/posts/feed", enableCORS(handler.GetFeed))
	http.HandleFunc("/api/posts/update", enableCORS(handler.UpdatePost))
	http.HandleFunc("/api/posts/delete", enableCORS(handler.DeletePost))

	// Comment routes
	http.HandleFunc("/api/comments", enableCORS(handler.CreateComment))
	http.HandleFunc("/api/comments/get", enableCORS(handler.GetComments))
	http.HandleFunc("/api/comments/react", enableCORS(handler.ReactToComment))

	// Post reaction routes
	http.HandleFunc("/api/posts/react", enableCORS(handler.ReactToPost))

	// Upload route
	http.HandleFunc("/api/upload", enableCORS(handler.UploadImage))

	// Group routes
	http.HandleFunc("/api/groups/create", enableCORS(handler.CreateGroup))
	http.HandleFunc("/api/groups/get", enableCORS(handler.GetGroup))
	http.HandleFunc("/api/groups/user", enableCORS(handler.GetUserGroups))
	http.HandleFunc("/api/groups/all", enableCORS(handler.GetAllGroups))
	http.HandleFunc("/api/groups/invite", enableCORS(handler.InviteToGroup))
	http.HandleFunc("/api/groups/respond", enableCORS(handler.RespondToGroupInvite))
	http.HandleFunc("/api/groups/join/request", enableCORS(handler.RequestToJoinGroup))
	http.HandleFunc("/api/groups/join/requests", enableCORS(handler.GetGroupJoinRequests))
	http.HandleFunc("/api/groups/join/respond", enableCORS(handler.RespondToJoinRequest))
	http.HandleFunc("/api/groups/posts/create", enableCORS(handler.CreateGroupPost))
	http.HandleFunc("/api/groups/posts/get", enableCORS(handler.GetGroupPosts))

	// Event routes
	http.HandleFunc("/api/events/create", enableCORS(handler.CreateEvent))
	http.HandleFunc("/api/events/get", enableCORS(handler.GetGroupEvents))
	http.HandleFunc("/api/events/respond", enableCORS(handler.RespondToEvent))

	// Message routes
	http.HandleFunc("/api/messages/private", enableCORS(handler.GetPrivateMessages))
	http.HandleFunc("/api/messages/group", enableCORS(handler.GetGroupMessages))
	http.HandleFunc("/api/messages/conversations", enableCORS(handler.GetConversations))

	// Notification routes
	http.HandleFunc("/api/notifications", enableCORS(handler.GetNotifications))
	http.HandleFunc("/api/notifications/read", enableCORS(handler.MarkNotificationRead))
	http.HandleFunc("/api/notifications/unread", enableCORS(handler.GetUnreadCount))

	// WebSocket route
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, repo, w, r)
	})

	log.Println("Server starting on :8080")
	log.Println("WebSocket available at ws://localhost:8080/ws")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}