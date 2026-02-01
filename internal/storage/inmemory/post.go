package inmemory

import (
	"PostService/internal/domain"
	"errors"
	"sync"
	"time"
)

type PostStorage struct {
	mu     sync.RWMutex
	posts  []domain.Post
	nextID int
}

func NewPostStorage() *PostStorage {
	return &PostStorage{
		posts:  make([]domain.Post, 0, 10),
		nextID: 1,
	}
}

func (s *PostStorage) GetAllPosts(limit, offset int) ([]domain.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if offset >= len(s.posts) {
		return []domain.Post{}, nil
	}

	end := offset + limit
	if end > len(s.posts) || limit <= 0 {
		end = len(s.posts)
	}

	result := make([]domain.Post, end-offset)
	copy(result, s.posts[offset:end])

	// переворачиваем список постов, чтобы сверху были новые
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

func (s *PostStorage) GetPostByID(id int) (*domain.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, p := range s.posts {
		if p.ID == id {
			pCopy := p
			return &pCopy, nil
		}
	}
	return nil, errors.New("post not found")
}

func (s *PostStorage) CreatePost(p domain.Post) (domain.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.ID = s.nextID
	p.CreatedAt = time.Now()
	s.nextID++

	s.posts = append(s.posts, p)
	return p, nil
}
