package service

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
)

type likeService struct {
	di             *internal.Di
	contextService domain.ContextService
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

	return &likeService{
		di:             di,
		likeRepository: likeRepository,
		contextService: contextService,
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
