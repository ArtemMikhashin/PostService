// internal/storage/postgres/comment.go
package postgres

import (
	"PostService/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CommentStorage struct {
	db *sqlx.DB
}

func NewCommentStorage(db *sqlx.DB) *CommentStorage {
	return &CommentStorage{db: db}
}

func (s *CommentStorage) CreateComment(c domain.Comment) (domain.Comment, error) {
	const query = `
		INSERT INTO comments (author, content, post_id, parent_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := s.db.QueryRow(
		query, c.Author, c.Content, c.PostID, c.ParentID,
	).Scan(&c.ID, &c.CreatedAt)

	return c, err
}

// where parentID null
func (s *CommentStorage) GetCommentsByPost(postID, limit, offset int) ([]domain.Comment, error) {
	const query = `
		SELECT id, created_at, author, content, post_id, parent_id
		FROM comments
		WHERE post_id = $1 AND parent_id IS NULL
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`

	var comments []domain.Comment
	err := s.db.Select(&comments, query, postID, limit, offset)
	return comments, err
}
