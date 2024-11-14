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

func TestGetFeed_WhenSuccess_ShouldReturnPostsWithLikes(t *testing.T) {
	ctx := context.Background()
	postServiceMock := new(mocks.PostService)
	likeServiceMock := new(mocks.LikeService)
	contextServiceMock := new(mocks.ContextService)

	feedService := &feedService{
		postService:    postServiceMock,
		likeService:    likeServiceMock,
		contextService: contextServiceMock,
	}

	page, limit := 1, 10
	userID := uuid.New()
	postID := uuid.New()
	posts := &domain.Pagination[*domain.PostResponse]{
		Rows: []*domain.PostResponse{{ID: postID}},
	}
	likes := map[uuid.UUID]bool{postID: true}

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	postServiceMock.On("GetPosts", ctx, page, limit).Return(posts, nil)
	likeServiceMock.On("UserLikedPosts", ctx, userID, mock.Anything).Return(likes, nil)

	result, err := feedService.GetFeed(ctx, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(posts.Rows), len(result.Rows))
	assert.True(t, result.Rows[0].LikesByUser)
	postServiceMock.AssertExpectations(t)
	likeServiceMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestGetFeed_WhenPostServiceFails_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	postServiceMock := new(mocks.PostService)
	likeServiceMock := new(mocks.LikeService)
	contextServiceMock := new(mocks.ContextService)

	feedService := &feedService{
		postService:    postServiceMock,
		likeService:    likeServiceMock,
		contextService: contextServiceMock,
	}

	page, limit := 1, 10
	postServiceMock.On("GetPosts", ctx, page, limit).Return(nil, errors.New("post service error"))

	result, err := feedService.GetFeed(ctx, page, limit)

	assert.ErrorContains(t, err, "post service error")
	assert.Nil(t, result)
	postServiceMock.AssertExpectations(t)
}

func TestGetFeed_WhenNoPostsFound_ShouldReturnErrPostNotFound(t *testing.T) {
	ctx := context.Background()
	postServiceMock := new(mocks.PostService)
	likeServiceMock := new(mocks.LikeService)
	contextServiceMock := new(mocks.ContextService)

	feedService := &feedService{
		postService:    postServiceMock,
		likeService:    likeServiceMock,
		contextService: contextServiceMock,
	}

	page, limit := 1, 10
	emptyPosts := &domain.Pagination[*domain.PostResponse]{Rows: []*domain.PostResponse{}}

	postServiceMock.On("GetPosts", ctx, page, limit).Return(emptyPosts, nil)

	result, err := feedService.GetFeed(ctx, page, limit)

	assert.ErrorIs(t, err, domain.ErrPostNotFound)
	assert.Nil(t, result)
	postServiceMock.AssertExpectations(t)
}

func TestGetFeed_WhenLikeServiceFails_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	postServiceMock := new(mocks.PostService)
	likeServiceMock := new(mocks.LikeService)
	contextServiceMock := new(mocks.ContextService)

	feedService := &feedService{
		postService:    postServiceMock,
		likeService:    likeServiceMock,
		contextService: contextServiceMock,
	}

	page, limit := 1, 10
	userID := uuid.New()
	postID := uuid.New()
	posts := &domain.Pagination[*domain.PostResponse]{
		Rows: []*domain.PostResponse{{ID: postID}},
	}

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	postServiceMock.On("GetPosts", ctx, page, limit).Return(posts, nil)
	likeServiceMock.On("UserLikedPosts", ctx, userID, mock.Anything).Return(nil, errors.New("like service error"))

	result, err := feedService.GetFeed(ctx, page, limit)

	assert.ErrorContains(t, err, "like service error")
	assert.Nil(t, result)
	postServiceMock.AssertExpectations(t)
	likeServiceMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}
