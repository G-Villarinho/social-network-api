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

func TestCreateLike_WhenUserAlreadyLikedPost_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	likeRepoMock := new(mocks.LikeRepository)
	contextServiceMock := new(mocks.ContextService)
	cacheMock := new(mocks.MemoryCacheRepository)

	likeService := &likeService{
		likeRepository: likeRepoMock,
		contextService: contextServiceMock,
		memoryCache:    cacheMock,
	}

	payload := domain.LikePayload{
		UserID: uuid.New(),
		PostID: uuid.New(),
	}

	likeRepoMock.On("UserLikedPost", ctx, payload.PostID, payload.UserID).Return(true, nil)

	err := likeService.CreateLike(ctx, payload)

	assert.Equal(t, domain.ErrPostAlreadyLiked, err)
	likeRepoMock.AssertExpectations(t)
}

func TestCreateLike_WhenRepositoryFails_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	likeRepoMock := new(mocks.LikeRepository)
	contextServiceMock := new(mocks.ContextService)
	cacheMock := new(mocks.MemoryCacheRepository)

	likeService := &likeService{
		likeRepository: likeRepoMock,
		contextService: contextServiceMock,
		memoryCache:    cacheMock,
	}

	payload := domain.LikePayload{
		UserID: uuid.New(),
		PostID: uuid.New(),
	}

	likeRepoMock.On("UserLikedPost", ctx, payload.PostID, payload.UserID).Return(false, nil)
	likeRepoMock.On("CreateLike", ctx, mock.Anything).Return(errors.New("repository error"))

	err := likeService.CreateLike(ctx, payload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository error")
	likeRepoMock.AssertExpectations(t)
}

func TestCreateLike_WhenSuccess_ShouldCreateLike(t *testing.T) {
	ctx := context.Background()
	likeRepoMock := new(mocks.LikeRepository)
	contextServiceMock := new(mocks.ContextService)
	cacheMock := new(mocks.MemoryCacheRepository)

	likeService := &likeService{
		likeRepository: likeRepoMock,
		contextService: contextServiceMock,
		memoryCache:    cacheMock,
	}

	payload := domain.LikePayload{
		UserID: uuid.New(),
		PostID: uuid.New(),
	}

	likeRepoMock.On("UserLikedPost", ctx, payload.PostID, payload.UserID).Return(false, nil)
	likeRepoMock.On("CreateLike", ctx, mock.Anything).Return(nil)

	err := likeService.CreateLike(ctx, payload)

	assert.NoError(t, err)
	likeRepoMock.AssertExpectations(t)
}

func TestUserLikedPosts_WhenCacheFails_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	likeRepoMock := new(mocks.LikeRepository)
	contextServiceMock := new(mocks.ContextService)
	cacheMock := new(mocks.MemoryCacheRepository)

	likeService := &likeService{
		likeRepository: likeRepoMock,
		contextService: contextServiceMock,
		memoryCache:    cacheMock,
	}

	userID := uuid.New()
	postIDs := []uuid.UUID{uuid.New(), uuid.New()}

	cacheMock.On("GetCachedLikes", ctx, userID, postIDs).Return(nil, errors.New("cache error"))

	result, err := likeService.UserLikedPosts(ctx, userID, postIDs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cache error")
	assert.Nil(t, result)
	cacheMock.AssertExpectations(t)
}

func TestUserLikedPosts_WhenDatabaseFetchFails_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	likeRepoMock := new(mocks.LikeRepository)
	contextServiceMock := new(mocks.ContextService)
	cacheMock := new(mocks.MemoryCacheRepository)

	likeService := &likeService{
		likeRepository: likeRepoMock,
		contextService: contextServiceMock,
		memoryCache:    cacheMock,
	}

	userID := uuid.New()
	postIDs := []uuid.UUID{uuid.New(), uuid.New()}
	likeCache := &domain.LikeCache{
		CachedLikes:  []uuid.UUID{},
		MissingLikes: postIDs,
	}

	cacheMock.On("GetCachedLikes", ctx, userID, postIDs).Return(likeCache, nil)
	likeRepoMock.On("UserLikedPosts", ctx, userID, postIDs).Return(nil, errors.New("database error"))

	result, err := likeService.UserLikedPosts(ctx, userID, postIDs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.Nil(t, result)
	cacheMock.AssertExpectations(t)
	likeRepoMock.AssertExpectations(t)
}

func TestUserLikedPosts_WhenSuccess_ShouldReturnLikesMap(t *testing.T) {
	ctx := context.Background()
	likeRepoMock := new(mocks.LikeRepository)
	contextServiceMock := new(mocks.ContextService)
	cacheMock := new(mocks.MemoryCacheRepository)

	likeService := &likeService{
		likeRepository: likeRepoMock,
		contextService: contextServiceMock,
		memoryCache:    cacheMock,
	}

	userID := uuid.New()
	postIDs := []uuid.UUID{uuid.New(), uuid.New()}
	cachedLikes := []uuid.UUID{postIDs[0]}
	likeCache := &domain.LikeCache{
		CachedLikes:  cachedLikes,
		MissingLikes: []uuid.UUID{postIDs[1]},
	}

	cacheMock.On("GetCachedLikes", ctx, userID, postIDs).Return(likeCache, nil)
	likeRepoMock.On("UserLikedPosts", ctx, userID, likeCache.MissingLikes).Return([]uuid.UUID{postIDs[1]}, nil)
	cacheMock.On("SetPostLike", ctx, postIDs[1], userID).Return(nil)

	result, err := likeService.UserLikedPosts(ctx, userID, postIDs)

	assert.NoError(t, err)
	assert.Equal(t, map[uuid.UUID]bool{
		postIDs[0]: true,
		postIDs[1]: true,
	}, result)
	cacheMock.AssertExpectations(t)
	likeRepoMock.AssertExpectations(t)
}

func TestDeleteLike_WhenUserDidNotLikePost_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	likeRepoMock := new(mocks.LikeRepository)
	contextServiceMock := new(mocks.ContextService)
	cacheMock := new(mocks.MemoryCacheRepository)

	likeService := &likeService{
		likeRepository: likeRepoMock,
		contextService: contextServiceMock,
		memoryCache:    cacheMock,
	}

	payload := domain.LikePayload{
		UserID: uuid.New(),
		PostID: uuid.New(),
	}

	likeRepoMock.On("UserLikedPost", ctx, payload.PostID, payload.UserID).Return(false, nil)

	err := likeService.DeleteLike(ctx, payload)

	assert.Equal(t, domain.ErrPostNotLiked, err)
	likeRepoMock.AssertExpectations(t)
}

func TestDeleteLike_WhenRepositoryFails_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	likeRepoMock := new(mocks.LikeRepository)
	contextServiceMock := new(mocks.ContextService)
	cacheMock := new(mocks.MemoryCacheRepository)

	likeService := &likeService{
		likeRepository: likeRepoMock,
		contextService: contextServiceMock,
		memoryCache:    cacheMock,
	}

	payload := domain.LikePayload{
		UserID: uuid.New(),
		PostID: uuid.New(),
	}

	likeRepoMock.On("UserLikedPost", ctx, payload.PostID, payload.UserID).Return(true, nil)
	likeRepoMock.On("DeleteLike", ctx, mock.Anything).Return(errors.New("repository error"))

	err := likeService.DeleteLike(ctx, payload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository error")
	likeRepoMock.AssertExpectations(t)
}

func TestDeleteLike_WhenSuccess_ShouldDeleteLike(t *testing.T) {
	ctx := context.Background()
	likeRepoMock := new(mocks.LikeRepository)
	contextServiceMock := new(mocks.ContextService)
	cacheMock := new(mocks.MemoryCacheRepository)

	likeService := &likeService{
		likeRepository: likeRepoMock,
		contextService: contextServiceMock,
		memoryCache:    cacheMock,
	}

	payload := domain.LikePayload{
		UserID: uuid.New(),
		PostID: uuid.New(),
	}

	likeRepoMock.On("UserLikedPost", ctx, payload.PostID, payload.UserID).Return(true, nil)
	likeRepoMock.On("DeleteLike", ctx, mock.Anything).Return(nil)

	err := likeService.DeleteLike(ctx, payload)

	assert.NoError(t, err)
	likeRepoMock.AssertExpectations(t)
}
