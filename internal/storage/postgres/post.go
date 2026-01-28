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

func (s *PostStorage) CreatePost(p domain.Post) (domain.Post, error) {
	const query = `
		INSERT INTO posts (author, title, content, comments_allowed)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := s.db.QueryRow(query,
		p.Author,
		p.Title,
		p.Content,
		p.CommentsAllowed,
	).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return domain.Post{}, err
	}
	return p, nil
}
