package graphql

import "PostService/internal/service"

type Resolver struct {
	PostService          *service.PostService
	CommentService       *service.CommentService
	CommentSubscriptions *CommentSubscriptions
}
