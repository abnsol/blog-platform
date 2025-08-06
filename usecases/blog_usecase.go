package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/blog-platform/domain"
)

type blogUsecase struct {
	blogRepo domain.IBlogRepository
}

func NewBlogUsecase(repo domain.IBlogRepository) domain.IBlogUsecase {
	return &blogUsecase{
		blogRepo: repo,
	}
}

func (uc *blogUsecase) CreateBlog(ctx context.Context, blog *domain.Blog, tags []string) error {
	if blog.Title == "" || blog.Content == "" {
		return errors.New("title and content cannot be empty")
	}

	err := uc.blogRepo.Create(ctx, blog)
	if err != nil {
		return fmt.Errorf("failed to create blog: %w", err)
	}

	// Ensure blog ID is populated after creation
	if blog.ID == 0 {
		return errors.New("blog ID not set after creation")
	}

	for _, tag := range tags {
		if tag == "" {
			continue // skip empty tags
		}

		tagID, err := uc.blogRepo.FindOrCreateTag(ctx, tag)
		if err != nil {
			return fmt.Errorf("failed to find or create tag '%s': %w", tag, err)
		}

		err = uc.blogRepo.LinkTagToBlog(ctx, int64(blog.ID), tagID)
		if err != nil {
			return fmt.Errorf("failed to link tag '%s' to blog: %w", tag, err)
		}
	}

	return nil
}

func (uc *blogUsecase) FetchAllBlogs(ctx context.Context) ([]*domain.Blog, error) {
	blogs, err := uc.blogRepo.FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blogs: %w", err)
	}
	return blogs, nil
}
