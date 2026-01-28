package inmemory

import (
	"PostService/internal/domain"
	"sync"
	"time"
)

type CommentStorage struct {
	mu        sync.RWMutex
	comments  []domain.Comment
	nextID    int
	postStore interface {
		GetAllPosts(limit, offset int) ([]domain.Post, error)
	}
}

func NewCommentStorage(postStore interface {
	GetAllPosts(limit, offset int) ([]domain.Post, error)
}) *CommentStorage {
	return &CommentStorage{
		comments:  make([]domain.Comment, 0, 10),
		nextID:    1,
		postStore: postStore,
	}
}

func (s *CommentStorage) CreateComment(c domain.Comment) (domain.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c.ID = s.nextID
	c.CreatedAt = time.Now()
	s.nextID++
	s.comments = append(s.comments, c)
	return c, nil
}

// Только главные комментарии (ParentID == nil)
func (s *CommentStorage) GetCommentsByPost(postID, limit, offset int) ([]domain.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []domain.Comment
	for _, comment := range s.comments {
		if comment.PostID == postID && comment.ParentID == nil {
			result = append(result, comment)
		}
	}

	if offset >= len(result) {
		return []domain.Comment{}, nil
	}
	end := offset + limit
	if end > len(result) || limit <= 0 {
		end = len(result)
	}
	return result[offset:end], nil
}

func (s *CommentStorage) GetReplies(parentID int) ([]domain.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []domain.Comment
	for _, c := range s.comments {
		if c.ParentID != nil && *c.ParentID == parentID {
			result = append(result, c)
		}
	}
	return result, nil
}
