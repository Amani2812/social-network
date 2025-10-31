package main

import (
	"log"
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/handlers"
	"social-network/pkg/models"
)

func main() {
	db, err := sqlite.NewDB("./social_network.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.RunMigrations("./pkg/db/migrations/sqlite"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := models.NewUserRepository(db.DB)
	followRepo := models.NewFollowRepository(db.DB)
	postRepo := models.NewPostRepository(db.DB)
	commentRepo := models.NewCommentRepository(db.DB)

	authHandler := handlers.NewAuthHandler(userRepo)
	followHandler := handlers.NewFollowHandler(followRepo, userRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	postHandler := handlers.NewPostHandler(postRepo, commentRepo, userRepo)
	commentHandler := handlers.NewCommentHandler(commentRepo, postRepo)

	// Routes
	http.HandleFunc("/api/auth/register", authHandler.Register)
	http.HandleFunc("/api/auth/login", authHandler.Login)
	http.HandleFunc("/api/auth/logout", authHandler.Logout)
	http.HandleFunc("/api/auth/me", authHandler.GetCurrentUser)

	// Follow routes
	http.HandleFunc("/api/follow/request", followHandler.SendFollowRequest)
	http.HandleFunc("/api/follow/respond", followHandler.RespondToFollowRequest)
	http.HandleFunc("/api/follow/unfollow", followHandler.Unfollow)
	http.HandleFunc("/api/follow/followers", followHandler.GetFollowers)
	http.HandleFunc("/api/follow/following", followHandler.GetFollowing)
	http.HandleFunc("/api/follow/pending", followHandler.GetPendingRequests)

	// User routes
	http.HandleFunc("/api/users/", userHandler.GetUserByID)
	http.HandleFunc("/api/users/search", userHandler.SearchUsers)
	http.HandleFunc("/api/users/profile", userHandler.UpdateProfile)

	// Post routes
	http.HandleFunc("/api/posts", postHandler.CreatePost)
	http.HandleFunc("/api/posts/get", postHandler.GetPost)
	http.HandleFunc("/api/posts/user", postHandler.GetUserPosts)
	http.HandleFunc("/api/posts/feed", postHandler.GetFeed)
	http.HandleFunc("/api/posts/update", postHandler.UpdatePost)
	http.HandleFunc("/api/posts/delete", postHandler.DeletePost)

	// Comment routes
	http.HandleFunc("/api/comments", commentHandler.CreateComment)
	http.HandleFunc("/api/comments/get", commentHandler.GetComments)
	http.HandleFunc("/api/comments/update", commentHandler.UpdateComment)
	http.HandleFunc("/api/comments/delete", commentHandler.DeleteComment)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
