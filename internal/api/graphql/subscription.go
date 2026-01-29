package graphql

import (
	"PostService/internal/domain"
	"sync"
)

type CommentSubscription struct {
	ch chan *domain.Comment
	id int
}

type CommentSubscriptions struct {
	mu      sync.Mutex
	chans   map[int][]CommentSubscription
	counter int
}

func NewCommentSubscriptions() *CommentSubscriptions {
	return &CommentSubscriptions{
		chans:   make(map[int][]CommentSubscription),
		counter: 0,
	}
}

func (c *CommentSubscriptions) CreateSubscription(postId int) (int, chan *domain.Comment, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan *domain.Comment, 10)
	c.counter++
	sub := CommentSubscription{ch: ch, id: c.counter}
	c.chans[postId] = append(c.chans[postId], sub)

	return c.counter, ch, nil
}

func (c *CommentSubscriptions) DeleteSubscription(postId, subId int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	subs, exists := c.chans[postId]
	if !exists {
		return nil
	}

	for i, sub := range subs {
		if sub.id == subId {
			close(sub.ch)
			c.chans[postId] = append(subs[:i], subs[i+1:]...)
			break
		}
	}

	return nil
}

func (c *CommentSubscriptions) Notify(postId int, comment domain.Comment) {
	c.mu.Lock()
	defer c.mu.Unlock()

	subs, exists := c.chans[postId]
	if !exists {
		return
	}

	for _, sub := range subs {
		select {
		case sub.ch <- &comment:
		default:
		}
	}
}
