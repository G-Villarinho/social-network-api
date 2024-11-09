package service

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
)

type feedService struct {
	di          *internal.Di
	postService domain.PostService
	likeService domain.LikeService
}

func NewFeedService(di *internal.Di) (domain.FeedService, error) {
	postService, err := internal.Invoke[domain.PostService](di)
	if err != nil {
		return nil, err
	}

	likeService, err := internal.Invoke[domain.LikeService](di)
	if err != nil {
		return nil, err
	}

	return &feedService{
		di:          di,
		postService: postService,
		likeService: likeService,
	}, nil
}

func (f *feedService) GenerateFeed(ctx context.Context, page int, limit int) (*domain.Pagination[*domain.PostResponse], error) {
	paginatedPosts, err := f.postService.GetPosts(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	return paginatedPosts, nil
}
