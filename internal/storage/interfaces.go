package storage

import "PostService/internal/domain"

type PostStorage interface {
	GetAllPosts(limit, offset int) ([]domain.Post, error)
	GetPostByID(id int) (*domain.Post, error)
	CreatePost(p domain.Post) (domain.Post, error)
}

type CommentStorage interface {
	CreateComment(c domain.Comment) (domain.Comment, error)
	GetCommentsByPost(postID, limit, offset int) ([]domain.Comment, error)
	GetReplies(parentID int) ([]domain.Comment, error)
}
