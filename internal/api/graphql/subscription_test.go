package graphql

import (
	"PostService/internal/domain"
	"testing"
	"time"
)

func TestCommentSubscriptions(t *testing.T) {
	subscriptions := NewCommentSubscriptions()

	t.Run("create and notify subscription", func(t *testing.T) {
		postID := 1
		id, ch, err := subscriptions.CreateSubscription(postID)
		if err != nil {
			t.Fatalf("failed to create subscription: %v", err)
		}

		comment := domain.Comment{
			ID:      1,
			Author:  "TestAuthor",
			Content: "TestContent",
			PostID:  postID,
		}

		subscriptions.Notify(postID, comment)

		select {
		case received := <-ch:
			if received.ID != comment.ID {
				t.Errorf("wrong comment ID: got %d, want %d", received.ID, comment.ID)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("timeout waiting for subscription notification")
		}

		err = subscriptions.DeleteSubscription(postID, id)
		if err != nil {
			t.Errorf("failed to delete subscription: %v", err)
		}
	})
}
