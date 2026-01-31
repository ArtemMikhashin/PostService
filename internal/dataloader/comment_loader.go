package dataloader

import (
	"PostService/internal/domain"
	"PostService/internal/service"
	"context"
	"github.com/graph-gophers/dataloader/v7"
	"net/http"
)

type commentLoaderKey struct{}

func CommentBatcher(commentService *service.CommentService) dataloader.BatchFunc[int, []*domain.Comment] {
	return func(ctx context.Context, parentIDs []int) []*dataloader.Result[[]*domain.Comment] {
		results := make([]*dataloader.Result[[]*domain.Comment], len(parentIDs))

		for i, parentID := range parentIDs {
			replies, err := commentService.GetReplies(parentID)
			if err != nil {
				results[i] = &dataloader.Result[[]*domain.Comment]{Error: err}
				continue
			}
			ptrReplies := make([]*domain.Comment, len(replies))
			for j := range replies {
				c := replies[j]
				ptrReplies[j] = &c
			}
			results[i] = &dataloader.Result[[]*domain.Comment]{Data: ptrReplies}
		}

		return results
	}
}

func LoaderFromContext(ctx context.Context) *dataloader.Loader[int, []*domain.Comment] {
	loader := ctx.Value(commentLoaderKey{}).(*dataloader.Loader[int, []*domain.Comment])
	return loader
}

func CommentLoaderMiddleware(commentService *service.CommentService) func(httpHandler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loader := dataloader.NewBatchedLoader(CommentBatcher(commentService))
			ctx := context.WithValue(r.Context(), commentLoaderKey{}, loader)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
