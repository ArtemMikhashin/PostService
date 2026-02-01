package service

import (
	"PostService/internal/consts"
	"PostService/internal/domain"
	"PostService/internal/storage"
	"PostService/internal/storage/inmemory"
	"PostService/internal/storage/postgres"
	"errors"
	"fmt"
)

type CommentService struct {
	storage    storage.CommentStorage
	postGetter *PostService
}

func NewCommentService(
	inMemory bool,
	postgresComment *postgres.CommentStorage,
	inmemoryComment *inmemory.CommentStorage,
	postGetter *PostService,
) *CommentService {
	var s storage.CommentStorage
	if inMemory {
		s = inmemoryComment
	} else {
		s = postgresComment
	}
	return &CommentService{
		storage:    s,
		postGetter: postGetter,
	}
}

func (s *CommentService) CreateComment(input domain.CreateCommentInput) (domain.Comment, error) {
	if input.Author == "" {
		return domain.Comment{}, errors.New("author is required")
	}
	if len(input.Content) > consts.MaxCommentLength {
		return domain.Comment{}, fmt.Errorf("comment length exceeds maximum of %d characters", consts.MaxCommentLength)
	}
	if input.Post <= 0 {
		return domain.Comment{}, errors.New("invalid post ID")
	}
	post, err := s.postGetter.GetPostByID(input.Post)
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
	limit, offset := consts.DefaultPageSizeComments, 0
	// TODO: переписать на курсорную пагинацию
	if page != nil && *page > 0 {
		offset = (*page - 1) * limit
	}
	if pageSize != nil && *pageSize > 0 {
		if *pageSize > consts.MaxPageSizeComments {
			limit = consts.MaxPageSizeComments
		} else {
			limit = *pageSize
		}
	}
	return s.storage.GetCommentsByPost(postID, limit, offset)
}

func (s *CommentService) GetReplies(parentID int) ([]domain.Comment, error) {
	return s.storage.GetReplies(parentID)
}
