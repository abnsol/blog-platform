package usecases

import (
	"context"
	"errors"

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

func (uc blogUsecase) CreateBlog(ctx context.Context, blog *domain.Blog, tags []string) error {
	err := uc.blogRepo.Create(ctx, blog)
	if err != nil {
		return errors.New("failed to create blog")
	}
	// prevent empty strings from being added
	if blog.Title == "" || blog.Content == "" {
		return errors.New("title and content cannot be empty")
	}

	for _, tag := range tags {
		tagID, err := uc.blogRepo.FindOrCreateTag(ctx, tag)
		if err != nil {
			return err
		}
		err = uc.blogRepo.LinkTagToBlog(ctx, int64(blog.ID), tagID)
		if err != nil {
			return err
		}
	}

	return nil
}
