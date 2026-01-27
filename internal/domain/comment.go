package domain

import "time"

type Comment struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Author    string    `db:"author" json:"author"`
	Content   string    `db:"content" json:"content"`
	PostID    int       `db:"post_id" json:"post_id"`
	ParentID  *int      `db:"parent_id" json:"parent_id"`
}
