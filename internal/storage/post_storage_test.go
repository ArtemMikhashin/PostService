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

func testPostStorage(t *testing.T, store storage.PostStorage) {
	t.Run("create and get post", func(t *testing.T) {
		input := domain.Post{
			Author:          "TestAuthor",
			Title:           "TestTitle",
			Content:         "TestContent",
			CommentsAllowed: true,
		}
		created, err := store.CreatePost(input)
		if err != nil {
			t.Fatalf("CreatePost failed: %v", err)
		}
		if created.ID == 0 {
			t.Error("post ID should not be zero")
		}
		if created.Author != input.Author {
			t.Errorf("expected author %s, got %s", input.Author, created.Author)
		}

		got, err := store.GetPostByID(created.ID)
		if err != nil {
			t.Fatalf("GetPostByID failed: %v", err)
		}
		if got.Title != created.Title {
			t.Errorf("title mismatch: expected %s, got %s", created.Title, got.Title)
		}
	})

	t.Run("get non-existent post", func(t *testing.T) {
		_, err := store.GetPostByID(999999)
		if err == nil {
			t.Error("expected error for non-existent post")
		}
	})

	t.Run("pagination", func(t *testing.T) {
		// Создаём 5 постов
		for i := 0; i < 5; i++ {
			_, err := store.CreatePost(domain.Post{
				Author:  "Author",
				Title:   "Title",
				Content: "Content",
			})
			if err != nil {
				t.Fatalf("failed to create post #%d: %v", i, err)
			}
		}

		posts, err := store.GetAllPosts(2, 1) // limit=2, offset=1 → 2-й и 3-й посты
		if err != nil {
			t.Fatalf("GetAllPosts failed: %v", err)
		}
		if len(posts) != 2 {
			t.Errorf("expected 2 posts, got %d", len(posts))
		}
	})
}

func TestInMemoryPostStorage(t *testing.T) {
	store := inmemory.NewPostStorage()
	testPostStorage(t, store)
}

func TestPostgresPostStorage(t *testing.T) {
	// Пропускаем тест, если нет подключения к БД (удобно для CI/локального запуска)
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

	// Очищаем таблицы перед тестом
	db.MustExec("TRUNCATE TABLE posts RESTART IDENTITY CASCADE")

	store := postgres.NewPostStorage(db)
	testPostStorage(t, store)
}
