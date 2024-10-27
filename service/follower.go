package service

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/google/uuid"
)

type followerService struct {
	di                 *pkg.Di
	followerRepository domain.FollowerRepository
	userRepository     domain.UserRepository
}

func NewFollowerService(di *pkg.Di) (domain.FollowerService, error) {
	followerRepository, err := pkg.Invoke[domain.FollowerRepository](di)
	if err != nil {
		return nil, err
	}

	userRepository, err := pkg.Invoke[domain.UserRepository](di)
	if err != nil {
		return nil, err
	}

	return &followerService{
		di:                 di,
		followerRepository: followerRepository,
		userRepository:     userRepository,
	}, nil
}

func (f *followerService) FollowUser(ctx context.Context, followerId uuid.UUID) error {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return domain.ErrSessionNotFound
	}

	if session.UserID == followerId {
		return domain.ErrUserCannotFollowItself
	}

	follower, err := f.userRepository.GetUserByID(ctx, followerId)
	if err != nil {
		return fmt.Errorf("error to get follower by ID: %w", err)
	}

	if follower == nil {
		return domain.ErrFollowerNotFound
	}

	following := &domain.Follower{
		UserID:     session.UserID,
		FollowerID: followerId,
	}

	if err := f.followerRepository.CreateFollower(ctx, *following); err != nil {
		return fmt.Errorf("error to create follower: %w", err)
	}

	return nil
}
