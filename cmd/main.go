package main

import (
	"PostService/internal/api/graphql"
	"PostService/internal/app"
	"PostService/internal/domain"
	"PostService/internal/storage/inmemory"
	"PostService/internal/storage/postgres"
	"PostService/pkg/logger"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	_ = godotenv.Load()
	log := logger.New()

	var postStore interface {
		GetAllPosts(limit, offset int) ([]domain.Post, error)
	}

	inMemory := os.Getenv("IN_MEMORY") == "true"
	if inMemory {
		log.Info.Println("in-memory storage")
		ps := inmemory.NewPostStorage()
		ps.AddSampleData()
		postStore = ps
	} else {
		log.Info.Println("PostgreSQL storage")
		db, err := app.ConnectDB()
		if err != nil {
			log.Error.Fatalf("DB connection failed: %v", err)
		}
		defer db.Close()
		postStore = postgres.NewPostStorage(db)
	}

	resolver := &graphql.Resolver{
		PostStore: postStore,
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))
	//TODO: разобраться с устаревшим методом

	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info.Printf("http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Error.Fatalf("Server failed: %v", err)
	}
}
