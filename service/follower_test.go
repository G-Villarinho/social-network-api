package service

import (
	"context"
	"errors"
	"testing"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFollowUser_Success(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	userID := uuid.New()
	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	userRepoMock.On("GetUserByID", ctx, userID).Return(&domain.User{ID: userID}, nil)
	followerRepoMock.On("GetFollower", ctx, userID, session.UserID).Return(nil, nil)
	followerRepoMock.On("CreateFollower", ctx, mock.Anything).Return(nil)

	err := followerService.FollowUser(ctx, userID)

	assert.NoError(t, err)
	userRepoMock.AssertExpectations(t)
	followerRepoMock.AssertExpectations(t)
}

func TestFollowUser_UserCannotFollowItself_ReturnsError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	err := followerService.FollowUser(ctx, session.UserID)

	assert.Equal(t, domain.ErrUserCannotFollowItself, err)
}

func TestFollowUser_UserNotFound_ReturnsError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	userID := uuid.New()
	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	userRepoMock.On("GetUserByID", ctx, userID).Return(nil, nil)

	err := followerService.FollowUser(ctx, userID)

	assert.Equal(t, domain.ErrFollowerNotFound, err)
	userRepoMock.AssertExpectations(t)
}

func TestFollowUser_SessionNotFound_ReturnsError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	userID := uuid.New()

	err := followerService.FollowUser(ctx, userID)

	assert.Equal(t, domain.ErrSessionNotFound, err)
}

func TestFollowUser_FollowerAlreadyExists_ReturnsError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	userID := uuid.New()
	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	userRepoMock.On("GetUserByID", ctx, userID).Return(&domain.User{ID: userID}, nil)
	followerRepoMock.On("GetFollower", ctx, userID, session.UserID).Return(&domain.Follower{UserID: userID, FollowerID: session.UserID}, nil)

	err := followerService.FollowUser(ctx, userID)

	assert.Equal(t, domain.ErrFollowerAlreadyExists, err)
	userRepoMock.AssertExpectations(t)
	followerRepoMock.AssertExpectations(t)
}

func TestFollowUser_CreateFollowerFails_ReturnsError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	userID := uuid.New()
	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	userRepoMock.On("GetUserByID", ctx, userID).Return(&domain.User{ID: userID}, nil)
	followerRepoMock.On("GetFollower", ctx, userID, session.UserID).Return(nil, nil)
	followerRepoMock.On("CreateFollower", ctx, mock.Anything).Return(errors.New("database error"))

	err := followerService.FollowUser(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error to create follower")
	userRepoMock.AssertExpectations(t)
	followerRepoMock.AssertExpectations(t)
}

func TestUnfollowUser_SessionNotFound_ReturnsError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	userID := uuid.New()

	err := followerService.UnfollowUser(ctx, userID)

	assert.Equal(t, domain.ErrSessionNotFound, err)
}

func TestUnfollowUser_UserCannotUnfollowItself_ReturnsError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	err := followerService.UnfollowUser(ctx, session.UserID)

	assert.Equal(t, domain.ErrUserCannotUnfollowItself, err)
}

func TestUnfollowUser_UserNotFound_ReturnsError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	userID := uuid.New()
	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	userRepoMock.On("GetUserByID", ctx, userID).Return(nil, nil)

	err := followerService.UnfollowUser(ctx, userID)

	assert.Equal(t, domain.ErrFollowerNotFound, err)
	userRepoMock.AssertExpectations(t)
}

func TestUnfollowUser_FollowerNotFound_ReturnsErrorFollowingNotFound(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	userID := uuid.New()
	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	userRepoMock.On("GetUserByID", ctx, userID).Return(&domain.User{ID: userID}, nil)
	followerRepoMock.On("GetFollower", ctx, userID, session.UserID).Return(nil, nil)

	err := followerService.UnfollowUser(ctx, userID)

	assert.Equal(t, domain.ErrFollowingNotFound, err)
	userRepoMock.AssertExpectations(t)
	followerRepoMock.AssertExpectations(t)
}

func TestUnfollowUser_DeleteFollowerFails_ReturnsError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	userID := uuid.New()
	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	userRepoMock.On("GetUserByID", ctx, userID).Return(&domain.User{ID: userID}, nil)
	followerRepoMock.On("GetFollower", ctx, userID, session.UserID).Return(&domain.Follower{ID: uuid.New(), UserID: userID, FollowerID: session.UserID}, nil)
	followerRepoMock.On("DeleteFollower", ctx, mock.Anything).Return(errors.New("database error"))

	err := followerService.UnfollowUser(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error to delete follower")
	userRepoMock.AssertExpectations(t)
	followerRepoMock.AssertExpectations(t)
}

func TestUnfollowUser_Success(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		userRepository:     userRepoMock,
		followerRepository: followerRepoMock,
	}

	userID := uuid.New()
	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	userRepoMock.On("GetUserByID", ctx, userID).Return(&domain.User{ID: userID}, nil)
	followerRepoMock.On("GetFollower", ctx, userID, session.UserID).Return(&domain.Follower{ID: uuid.New(), UserID: userID, FollowerID: session.UserID}, nil)
	followerRepoMock.On("DeleteFollower", ctx, mock.Anything).Return(nil)

	err := followerService.UnfollowUser(ctx, userID)

	assert.NoError(t, err)
	userRepoMock.AssertExpectations(t)
	followerRepoMock.AssertExpectations(t)
}

func TestGetFollowers_SessionNotFound_ReturnsError(t *testing.T) {
	ctx := context.Background()
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		followerRepository: followerRepoMock,
	}

	followers, err := followerService.GetFollowers(ctx)

	assert.Equal(t, domain.ErrSessionNotFound, err)
	assert.Nil(t, followers)
}

func TestGetFollowers_RepositoryError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		followerRepository: followerRepoMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	followerRepoMock.On("GetFollowers", ctx, session.UserID).Return(nil, errors.New("database error"))

	followers, err := followerService.GetFollowers(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error to get followers")
	assert.Nil(t, followers)
	followerRepoMock.AssertExpectations(t)
}

func TestGetFollowers_NoFollowersFound_ReturnsErrorFollowerNotFound(t *testing.T) {
	ctx := context.Background()
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		followerRepository: followerRepoMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	followerRepoMock.On("GetFollowers", ctx, session.UserID).Return(nil, nil)

	followers, err := followerService.GetFollowers(ctx)

	assert.Equal(t, domain.ErrFollowerNotFound, err)
	assert.Nil(t, followers)
	followerRepoMock.AssertExpectations(t)
}

func TestGetFollowers_Success_ReturnsFollowersResponse(t *testing.T) {
	ctx := context.Background()
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		followerRepository: followerRepoMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	follower := &domain.Follower{
		ID:         uuid.New(),
		UserID:     session.UserID,
		FollowerID: uuid.New(),
	}
	followers := []*domain.Follower{follower}

	followerRepoMock.On("GetFollowers", ctx, session.UserID).Return(followers, nil)

	followersResponse, err := followerService.GetFollowers(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, followersResponse)
	assert.Len(t, followersResponse, 1)
	assert.Equal(t, follower.ToFollowerResponse(), followersResponse[0])
	followerRepoMock.AssertExpectations(t)
}

func TestGetFollowings_SessionNotFound_ReturnsError(t *testing.T) {
	ctx := context.Background()
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		followerRepository: followerRepoMock,
	}

	following, err := followerService.GetFollowings(ctx)

	assert.Equal(t, domain.ErrSessionNotFound, err)
	assert.Nil(t, following)
}

func TestGetFollowings_RepositoryError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		followerRepository: followerRepoMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	followerRepoMock.On("GetFollowings", ctx, session.UserID).Return(nil, errors.New("database error"))

	following, err := followerService.GetFollowings(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error to get following")
	assert.Nil(t, following)
	followerRepoMock.AssertExpectations(t)
}

func TestGetFollowings_NoFollowingFound_ReturnsErrorFollowingNotFound(t *testing.T) {
	ctx := context.Background()
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		followerRepository: followerRepoMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	followerRepoMock.On("GetFollowings", ctx, session.UserID).Return(nil, nil)

	following, err := followerService.GetFollowings(ctx)

	assert.Equal(t, domain.ErrFollowingNotFound, err)
	assert.Nil(t, following)
	followerRepoMock.AssertExpectations(t)
}

func TestGetFollowings_Success_ReturnsFollowingResponse(t *testing.T) {
	ctx := context.Background()
	followerRepoMock := new(mocks.FollowerRepository)

	followerService := &followerService{
		followerRepository: followerRepoMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	following := []*domain.Follower{
		{
			ID:         uuid.New(),
			UserID:     session.UserID,
			FollowerID: uuid.New(),
		},
	}
	followerRepoMock.On("GetFollowings", ctx, session.UserID).Return(following, nil)

	followingResponse, err := followerService.GetFollowings(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, followingResponse)
	assert.Len(t, followingResponse, 1)
	assert.Equal(t, following[0].ToFollowerResponse(), followingResponse[0])
	followerRepoMock.AssertExpectations(t)
}