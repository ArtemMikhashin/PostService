package service

import (
	"PostService/internal/domain"
	"PostService/internal/storage/inmemory"
	"PostService/internal/storage/postgres"
	"errors"
)

type PostService struct {
	storage interface {
		GetAllPosts(limit, offset int) ([]domain.Post, error)
		CreatePost(p domain.Post) (domain.Post, error)
	}
}

func NewPostService(inMemory bool, postgresStorage *postgres.PostStorage, inmemoryStorage *inmemory.PostStorage) *PostService {
	var storage interface {
		GetAllPosts(limit, offset int) ([]domain.Post, error)
		CreatePost(p domain.Post) (domain.Post, error)
	}

	if inMemory {
		storage = inmemoryStorage
	} else {
		storage = postgresStorage
	}

	return &PostService{storage: storage}
}

func (s *PostService) CreatePost(input domain.CreatePostInput) (domain.Post, error) {
	if input.Author == "" {
		return domain.Post{}, errors.New("author is required")
	}

	if len(input.Content) > 2000 {
		return domain.Post{}, errors.New("content must be at most 2000 characters")
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
	limit, offset := 10, 0 //дефолтные
	if page != nil && *page > 0 {
		offset = (*page - 1) * limit
	}
	if pageSize != nil && *pageSize > 0 {
		limit = *pageSize
	}

	return s.storage.GetAllPosts(limit, offset)
}
