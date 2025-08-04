package domain

import (
	"context"
)

type IBlogRepository interface {
	Create(ctx context.Context, blog *Blog) error
	FindOrCreateTag(ctx context.Context, tagName string) (int64, error)
	LinkTagToBlog(ctx context.Context, blogID int64, tagID int64) error
	FetchAll(ctx context.Context) ([]*Blog, error)
}

type IBlogUsecase interface {
	CreateBlog(ctx context.Context, blog *Blog, tags []string) error
	FetchAllBlogs(ctx context.Context) ([]*Blog, error)
}
