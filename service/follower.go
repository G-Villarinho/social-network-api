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

	following, err := f.followerRepository.GetFollower(ctx, session.UserID, followerId)
	if err != nil {
		return fmt.Errorf("error to get follower: %w", err)
	}

	if following != nil {
		return domain.ErrFollowerAlreadyExists
	}

	following = &domain.Follower{
		UserID:     session.UserID,
		FollowerID: followerId,
	}

	if err := f.followerRepository.CreateFollower(ctx, *following); err != nil {
		return fmt.Errorf("error to create follower: %w", err)
	}

	return nil
}

func (f *followerService) UnfollowUser(ctx context.Context, followerId uuid.UUID) error {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return domain.ErrSessionNotFound
	}

	if session.UserID == followerId {
		return domain.ErrUserCannotUnfollowItself
	}

	follower, err := f.userRepository.GetUserByID(ctx, followerId)
	if err != nil {
		return fmt.Errorf("error to get follower by ID: %w", err)
	}

	if follower == nil {
		return domain.ErrFollowerNotFound
	}

	following, err := f.followerRepository.GetFollower(ctx, session.UserID, followerId)
	if err != nil {
		return fmt.Errorf("error to get follower: %w", err)
	}

	if following == nil {
		return domain.ErrFollowingNotFound
	}

	if err := f.followerRepository.DeleteFollower(ctx, followerId); err != nil {
		return fmt.Errorf("error to delete follower: %w", err)
	}

	return nil
}

func (f *followerService) GetFollowers(ctx context.Context) ([]*domain.FollowerResponse, error) {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return nil, domain.ErrSessionNotFound
	}

	followers, err := f.followerRepository.GetFollowers(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("error to get followers: %w", err)
	}

	if followers == nil {
		return nil, domain.ErrFollowerNotFound
	}

	var followersResponse []*domain.FollowerResponse
	for _, follower := range followers {
		followersResponse = append(followersResponse, follower.ToFollowerResponse())
	}

	return followersResponse, nil
}

func (f *followerService) GetFollowings(ctx context.Context) ([]*domain.FollowerResponse, error) {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return nil, domain.ErrSessionNotFound
	}

	following, err := f.followerRepository.GetFollowings(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("error to get following: %w", err)
	}

	if following == nil {
		return nil, domain.ErrFollowingNotFound
	}

	var followingResponse []*domain.FollowerResponse
	for _, follower := range following {
		followingResponse = append(followingResponse, follower.ToFollowerResponse())
	}

	return followingResponse, nil
}
