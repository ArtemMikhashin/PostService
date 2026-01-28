package main

import (
	"PostService/internal/api/graphql"
	"PostService/internal/app"
	"PostService/internal/service"
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

	inMemory := os.Getenv("IN_MEMORY") == "true"

	var PostPostgresStorage *postgres.PostStorage
	var PostInmemoryStorage *inmemory.PostStorage

	if !inMemory {
		db, err := app.ConnectDB()
		if err != nil {
			log.Error.Fatalf("DB connection failed: %v", err)
		}
		defer db.Close()
		PostPostgresStorage = postgres.NewPostStorage(db)
	} else {
		PostInmemoryStorage = inmemory.NewPostStorage()
	}

	postService := service.NewPostService(inMemory, PostPostgresStorage, PostInmemoryStorage)

	resolver := &graphql.Resolver{
		PostService: postService,
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

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
