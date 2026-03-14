package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"practice5/internal/db"
	"practice5/internal/handler"
	"practice5/internal/repository"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Wire up layers
	repo := repository.NewRepository(database)
	userHandler := handler.NewUserHandler(repo)

	// Routes
	mux := http.NewServeMux()

	// GET /users         - paginated + filtered + sorted list
	// GET /users/common-friends - common friends of 2 users
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUsers(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/common-friends", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		userHandler.GetCommonFriends(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
