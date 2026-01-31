package service

import (
	"PostService/internal/consts"
	"PostService/internal/domain"
	"PostService/internal/storage/inmemory"
	"testing"
)

func TestCommentService_CreateComment(t *testing.T) {
	postStorage := inmemory.NewPostStorage()
	commentStorage := inmemory.NewCommentStorage(postStorage)

	postService := NewPostService(true, nil, postStorage)
	commentService := NewCommentService(true, nil, commentStorage, postService)

	postInput := domain.CreatePostInput{
		Author:          "TestAuthor",
		Title:           "TestTitle",
		Content:         "TestContent",
		CommentsAllowed: false,
	}
	post, _ := postService.CreatePost(postInput)

	t.Run("comment on locked post", func(t *testing.T) {
		input := domain.CreateCommentInput{
			Author:  "TestAuthor",
			Content: "TestContent",
			Post:    post.ID,
		}
		_, err := commentService.CreateComment(input)
		if err == nil {
			t.Error("expected error commentsAllowed=false, got nil")
		}
	})

	t.Run("valid comment", func(t *testing.T) {
		allowedPost, _ := postService.CreatePost(domain.CreatePostInput{
			Author:          "TestAuthor",
			Title:           "TestTitle",
			Content:         "TestContent",
			CommentsAllowed: true,
		})

		input := domain.CreateCommentInput{
			Author:  "TestAuthor",
			Content: "TestContent",
			Post:    allowedPost.ID,
		}
		comment, err := commentService.CreateComment(input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if comment.ID != 1 {
			t.Errorf("expected ID=1 for comment, got %d", comment.ID)
		}
		if comment.PostID != allowedPost.ID {
			t.Errorf("wrong post ID in comment")
		}
	})

	t.Run("empty author", func(t *testing.T) {
		input := domain.CreateCommentInput{
			Author:  "",
			Content: "TestContent",
			Post:    1,
		}
		_, err := commentService.CreateComment(input)
		if err == nil {
			t.Error("expected error empty author")
		}
	})

	t.Run("content too long", func(t *testing.T) {
		input := domain.CreateCommentInput{
			Author:  "Spammer",
			Content: string(make([]byte, consts.MaxCommentLength+1)),
			Post:    1,
		}
		_, err := commentService.CreateComment(input)
		if err == nil {
			t.Errorf("expected error length content > %d", consts.MaxCommentLength)
		}
	})
}
