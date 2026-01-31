package service

import (
	"PostService/internal/consts"
	"PostService/internal/domain"
	"PostService/internal/storage/inmemory"
	"PostService/internal/storage/postgres"
	"errors"
	"fmt"
)

type PostService struct {
	storage interface {
		GetAllPosts(limit, offset int) ([]domain.Post, error)
		GetPostByID(id int) (*domain.Post, error)
		CreatePost(p domain.Post) (domain.Post, error)
	}
}

func NewPostService(inMemory bool, postgresStorage *postgres.PostStorage, inmemoryStorage *inmemory.PostStorage) *PostService {
	var storage interface {
		GetAllPosts(limit, offset int) ([]domain.Post, error)
		GetPostByID(id int) (*domain.Post, error)
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
