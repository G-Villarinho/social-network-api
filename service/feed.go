package service

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/google/uuid"
)

type feedService struct {
	di             *internal.Di
	postService    domain.PostService
	contextService domain.ContextService
	likeService    domain.LikeService
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

	contextService, err := internal.Invoke[domain.ContextService](di)
	if err != nil {
		return nil, err
	}

	return &feedService{
		di:             di,
		postService:    postService,
		likeService:    likeService,
		contextService: contextService,
	}, nil
}

func (f *feedService) GetFeed(ctx context.Context, page int, limit int) (*domain.Pagination[*domain.PostResponse], error) {
	paginatedPosts, err := f.postService.GetPosts(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	if len(paginatedPosts.Rows) == 0 {
		return nil, domain.ErrPostNotFound
	}

	postIDs := make([]uuid.UUID, len(paginatedPosts.Rows))
	for i, post := range paginatedPosts.Rows {
		postIDs[i] = post.ID
	}

	likes, err := f.likeService.UserLikedPosts(ctx, f.contextService.GetUserID(ctx), postIDs)
	if err != nil {
		return nil, fmt.Errorf("fetch user liked posts: %w", err)
	}

	for i, post := range paginatedPosts.Rows {
		liked, ok := likes[post.ID]
		if ok {
			paginatedPosts.Rows[i].SetLikesByUser(liked)
		}
	}

	return paginatedPosts, nil
}
