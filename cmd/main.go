package main

import (
	"PostService/internal/app"
	"PostService/internal/storage/postgres"
	"PostService/pkg/logger"
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

func main() {
	_ = godotenv.Load()
	log := logger.New()

	db, err := app.ConnectDB()
	if err != nil {
		log.Error.Fatalf("DB connection failed: %v", err)
	}
	defer db.Close()

	postStorage := postgres.NewPostStorage(db)

	http.HandleFunc("/debug/posts", func(w http.ResponseWriter, r *http.Request) {
		posts, err := postStorage.GetAllPosts(10, 0)
		if err != nil {
			log.Error.Printf("Failed to fetch posts: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			log.Error.Printf("Failed to encode response: %v", err)
		}
	})

	log.Info.Println("Connected to database")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			http.Error(w, "Database unreachable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Error.Fatalf("Server failed: %v", err)
	}
}
