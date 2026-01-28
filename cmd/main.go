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

	var postPgStore *postgres.PostStorage
	var postMemStore *inmemory.PostStorage
	var commentPgStore *postgres.CommentStorage
	var commentMemStore *inmemory.CommentStorage

	inMemory := os.Getenv("IN_MEMORY") == "true"
	if !inMemory {
		db, err := app.ConnectDB()
		if err != nil {
			log.Error.Fatalf("DB connection failed: %v", err)
		}
		defer db.Close()
		postPgStore = postgres.NewPostStorage(db)
		commentPgStore = postgres.NewCommentStorage(db)
	} else {
		postMemStore = inmemory.NewPostStorage()
		commentMemStore = inmemory.NewCommentStorage(postMemStore)
	}

	postService := service.NewPostService(inMemory, postPgStore, postMemStore)
	commentService := service.NewCommentService(inMemory, commentPgStore, commentMemStore, postService)

	resolver := &graphql.Resolver{
		PostService:    postService,
		CommentService: commentService,
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
