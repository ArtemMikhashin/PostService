package domain

import "time"

type Post struct {
	ID              int       `db:"id" json:"id"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	Author          string    `db:"author" json:"author"`
	Title           string    `db:"title" json:"title"`
	Content         string    `db:"content" json:"content"`
	CommentsAllowed bool      `db:"comments_allowed" json:"comments_allowed"`
}
