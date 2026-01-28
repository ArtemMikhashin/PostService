package service

import (
	"PostService/internal/domain"
	"PostService/internal/storage/inmemory"
	"PostService/internal/storage/postgres"
	"errors"
)

type CommentService struct {
	storage interface {
		CreateComment(c domain.Comment) (domain.Comment, error)
		GetCommentsByPost(postID, limit, offset int) ([]domain.Comment, error)
		GetReplies(parentID int) ([]domain.Comment, error)
	}
	postGetter interface {
		GetAllPosts(page, pageSize *int) ([]domain.Post, error)
	}
}

func NewCommentService(
	inMemory bool,
	postgresComment *postgres.CommentStorage,
	inmemoryComment *inmemory.CommentStorage,
	postGetter interface {
		GetAllPosts(page, pageSize *int) ([]domain.Post, error)
	},
) *CommentService {
	var storage interface {
		CreateComment(c domain.Comment) (domain.Comment, error)
		GetCommentsByPost(postID, limit, offset int) ([]domain.Comment, error)
		GetReplies(parentID int) ([]domain.Comment, error)
	}

	if inMemory {
		storage = inmemoryComment
	} else {
		storage = postgresComment
	}

	return &CommentService{
		storage:    storage,
		postGetter: postGetter,
	}
}

func (s *CommentService) CreateComment(input domain.CreateCommentInput) (domain.Comment, error) {
	if input.Author == "" {
		return domain.Comment{}, errors.New("author is required")
	}
	if len(input.Content) > 2000 {
		return domain.Comment{}, errors.New("content must be at most 2000 characters")
	}
	if input.Post <= 0 {
		return domain.Comment{}, errors.New("invalid post ID")
	}

	post, err := s.postGetter.(*PostService).GetPostByID(input.Post)
	if err != nil {
		return domain.Comment{}, errors.New("post not found")
	}

	if !post.CommentsAllowed {
		return domain.Comment{}, errors.New("comments are not allowed for this post")
	}

	comment := domain.Comment{
		Author:   input.Author,
		Content:  input.Content,
		PostID:   input.Post,
		ParentID: input.ReplyTo,
	}

	return s.storage.CreateComment(comment)
}

func (s *CommentService) GetCommentsByPost(postID int, page, pageSize *int) ([]domain.Comment, error) {
	limit, offset := 10, 0
	if page != nil && *page > 0 {
		offset = (*page - 1) * limit
	}
	if pageSize != nil && *pageSize > 0 {
		limit = *pageSize
	}
	return s.storage.GetCommentsByPost(postID, limit, offset)
}

func (s *CommentService) GetReplies(parentID int) ([]domain.Comment, error) {
	return s.storage.GetReplies(parentID)
}
