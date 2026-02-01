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

type PostService struct {
	storage storage.PostStorage
}

func NewPostService(inMemory bool, postgresStorage *postgres.PostStorage, inmemoryStorage *inmemory.PostStorage) *PostService {
	var store storage.PostStorage
	if inMemory {
		store = inmemoryStorage
	} else {
		store = postgresStorage
	}
	return &PostService{storage: store}
}

func (s *PostService) CreatePost(input domain.CreatePostInput) (domain.Post, error) {
	if input.Author == "" {
		return domain.Post{}, errors.New("author is required")
	}
	if len(input.Title) > consts.MaxPostTitleLength {
		return domain.Post{}, fmt.Errorf("title must be at most %d characters", consts.MaxPostTitleLength)
	}
	if len(input.Content) > consts.MaxPostContentLength {
		return domain.Post{}, fmt.Errorf("content must be at most %d characters", consts.MaxPostContentLength)
	}

	post := domain.Post{
		Author:          input.Author,
		Title:           input.Title,
		Content:         input.Content,
		CommentsAllowed: input.CommentsAllowed,
	}

	return s.storage.CreatePost(post)
}

func (s *PostService) GetAllPosts(page, pageSize *int) ([]domain.Post, error) {
	limit, offset := consts.DefaultPageSizePosts, 0

	if page != nil && *page > 0 {
		offset = (*page - 1) * limit
	}
	if pageSize != nil && *pageSize > 0 {
		if *pageSize > consts.MaxPageSizePosts {
			limit = consts.MaxPageSizePosts
		} else {
			limit = *pageSize
		}
	}

	return s.storage.GetAllPosts(limit, offset)
}

func (s *PostService) GetPostByID(id int) (*domain.Post, error) {
	return s.storage.GetPostByID(id)
}
