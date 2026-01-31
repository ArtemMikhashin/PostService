package service

import (
	"PostService/internal/consts"
	"PostService/internal/domain"
	"PostService/internal/storage/inmemory"
	"testing"
)

func TestPostService_CreatePost(t *testing.T) {
	storage := inmemory.NewPostStorage()
	service := NewPostService(true, nil, storage)

	t.Run("valid post", func(t *testing.T) {
		input := domain.CreatePostInput{
			Author:          "TestAuthor",
			Title:           "TestTitle",
			Content:         "TestContent",
			CommentsAllowed: true,
		}
		post, err := service.CreatePost(input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if post.ID != 1 {
			t.Errorf("expected ID=1 for post, got %d", post.ID)
		}
		if post.Author != "TestAuthor" {
			t.Errorf("expected author=TestAuthor, got %s", post.Author)
		}
	})

	t.Run("empty author", func(t *testing.T) {
		input := domain.CreatePostInput{
			Author:  "",
			Title:   "TestTitle",
			Content: "TestContent",
		}
		_, err := service.CreatePost(input)
		if err == nil {
			t.Error("expected error empty author")
		}
	})

	t.Run("content too long", func(t *testing.T) {
		input := domain.CreatePostInput{
			Author:  "TestAuthor",
			Title:   "TestTitle",
			Content: string(make([]byte, consts.MaxPostContentLength)),
		}
		_, err := service.CreatePost(input)
		if err == nil {
			t.Error("expected error length content")
		}
	})
}
