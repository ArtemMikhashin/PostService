package main

import (
	"PostService/internal/app"
	"PostService/pkg/logger"
	"context"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

func main() {
	_ = godotenv.Load()
	log := logger.New()

	// Подключаемся к БД
	db, err := app.ConnectDB()
	if err != nil {
		log.Error.Fatalf("DB connection failed: %v", err)
	}
	defer db.Close()

	log.Info.Println("Connected to database")

	// Health check
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
