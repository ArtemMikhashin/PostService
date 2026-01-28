package graphql

import "PostService/internal/domain"

type Resolver struct {
	PostStore interface {
		GetAllPosts(limit, offset int) ([]domain.Post, error)
	}
}
