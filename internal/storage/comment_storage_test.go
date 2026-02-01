package storage_test

import (
	"PostService/internal/domain"
	"PostService/internal/storage"
	"PostService/internal/storage/inmemory"
	"PostService/internal/storage/postgres"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func testCommentStorage(t *testing.T, commentStore storage.CommentStorage, postStore storage.PostStorage) {
	post, err := postStore.CreatePost(domain.Post{
		Author:          "Author",
		Title:           "Title",
		Content:         "Content",
		CommentsAllowed: true,
	})
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	t.Run("create root comment", func(t *testing.T) {
		comment, err := commentStore.CreateComment(domain.Comment{
			Author:   "Commenter",
			Content:  "Root comment",
			PostID:   post.ID,
			ParentID: nil,
		})
		if err != nil {
			t.Fatalf("CreateComment failed: %v", err)
		}
		if comment.ID == 0 {
			t.Error("comment ID should not be zero")
		}
	})

	t.Run("create reply", func(t *testing.T) {
		root, err := commentStore.CreateComment(domain.Comment{
			Author:   "Root",
			Content:  "Root",
			PostID:   post.ID,
			ParentID: nil,
		})
		if err != nil {
			t.Fatal(err)
		}

		reply, err := commentStore.CreateComment(domain.Comment{
			Author:   "Replier",
			Content:  "Reply",
			PostID:   post.ID,
			ParentID: &root.ID,
		})
		if err != nil {
			t.Fatalf("failed to create reply: %v", err)
		}

		replies, err := commentStore.GetReplies(root.ID)
		if err != nil {
			t.Fatalf("GetReplies failed: %v", err)
		}
		if len(replies) != 1 {
			t.Errorf("expected 1 reply, got %d", len(replies))
		}
		if replies[0].ID != reply.ID {
			t.Error("reply ID mismatch")
		}
	})

	t.Run("get comments by post (root only)", func(t *testing.T) {
		postStore := inmemory.NewPostStorage()
		commentStore := inmemory.NewCommentStorage(postStore)
		c1, _ := commentStore.CreateComment(domain.Comment{Author: "A", Content: "C1", PostID: post.ID, ParentID: nil})
		_, _ = commentStore.CreateComment(domain.Comment{Author: "B", Content: "C2", PostID: post.ID, ParentID: nil})
		_, _ = commentStore.CreateComment(domain.Comment{Author: "R", Content: "Reply", PostID: post.ID, ParentID: &c1.ID})

		comments, err := commentStore.GetCommentsByPost(post.ID, 10, 0)
		if err != nil {
			t.Fatalf("GetCommentsByPost failed: %v", err)
		}
		// Должны вернуться только корневые комментарии (c1, c2)
		if len(comments) != 2 {
			t.Errorf("expected 2 root comments, got %d", len(comments))
		}
	})
}

func TestInMemoryCommentStorage(t *testing.T) {
	postStore := inmemory.NewPostStorage()
	commentStore := inmemory.NewCommentStorage(postStore)
	testCommentStorage(t, commentStore, postStore)
}

func TestPostgresCommentStorage(t *testing.T) {
	if os.Getenv("TEST_PG") == "" {
		t.Skip("set TEST_PG=1 to run PostgreSQL tests")
	}

	_ = godotenv.Load()

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to test DB: %v", err)
	}
	defer db.Close()

	db.MustExec("TRUNCATE TABLE comments, posts RESTART IDENTITY CASCADE")

	postStore := postgres.NewPostStorage(db)
	commentStore := postgres.NewCommentStorage(db)
	testCommentStorage(t, commentStore, postStore)
}
