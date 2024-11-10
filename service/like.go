package service

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/google/uuid"
)

type likeService struct {
	di             *internal.Di
	contextService domain.ContextService
	memoryCache    domain.MemoryCacheRepository
	likeRepository domain.LikeRepository
}

func NewLikeService(di *internal.Di) (domain.LikeService, error) {
	likeRepository, err := internal.Invoke[domain.LikeRepository](di)
	if err != nil {
		return nil, err
	}

	contextService, err := internal.Invoke[domain.ContextService](di)
	if err != nil {
		return nil, err
	}

	memoryCache, err := internal.Invoke[domain.MemoryCacheRepository](di)
	if err != nil {
		return nil, err
	}

	return &likeService{
		di:             di,
		likeRepository: likeRepository,
		contextService: contextService,
		memoryCache:    memoryCache,
	}, nil
}

func (l *likeService) CreateLike(ctx context.Context, payload domain.LikePayload) error {
	userLiked, err := l.likeRepository.UserLikedPost(ctx, payload.PostID, payload.UserID)
	if err != nil {
		return fmt.Errorf("check if user liked post: %w", err)
	}

	if userLiked {
		return domain.ErrPostAlreadyLiked
	}

	like := payload.ToLike()

	if err := l.likeRepository.CreateLike(ctx, *like); err != nil {
		return err
	}

	return nil
}

func (l *likeService) UserLikedPosts(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	likeCache, err := l.memoryCache.GetCachedLikes(ctx, userID, postIDs)
	if err != nil {
		return nil, fmt.Errorf("error fetching likes from cache: %w", err)
	}

	likesMap := make(map[uuid.UUID]bool, len(postIDs))

	for _, likedPostID := range likeCache.CachedLikes {
		likesMap[likedPostID] = true
	}

	if len(likeCache.MissingLikes) > 0 {
		missingLikes, err := l.likeRepository.UserLikedPosts(ctx, userID, likeCache.MissingLikes)
		if err != nil {
			return nil, fmt.Errorf("error fetching missing likes from database: %w", err)
		}

		missingLikesMap := make(map[uuid.UUID]bool, len(missingLikes))
		for _, likedPostID := range missingLikes {
			missingLikesMap[likedPostID] = true
		}

		for _, postID := range likeCache.MissingLikes {
			liked := missingLikesMap[postID]
			likesMap[postID] = liked
			if err := l.memoryCache.SetPostLike(ctx, postID, userID); err != nil {
				return nil, fmt.Errorf("error setting like in cache: %w", err)
			}
		}
	}

	return likesMap, nil
}
