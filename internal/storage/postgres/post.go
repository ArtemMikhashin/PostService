package postgres

import (
	"PostService/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PostStorage struct {
	db *sqlx.DB
}

func NewPostStorage(db *sqlx.DB) *PostStorage {
	return &PostStorage{db: db}
}

func (s *PostStorage) GetAllPosts(limit, offset int) ([]domain.Post, error) {
	const query = `
		SELECT id, created_at, author, title, content, comments_allowed
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	var posts []domain.Post
	err := s.db.Select(&posts, query, limit, offset)
	return posts, err
}
