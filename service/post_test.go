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

func TestGetPosts_SuccessFromCache(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		postRepository:        postRepoMock,
		contextService:        contextServiceMock,
	}

	userID := uuid.New()
	page, limit := 1, 10
	cachedPosts := &domain.Pagination[*domain.PostResponse]{}

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("GetPosts", ctx, userID, page, limit).Return(cachedPosts, nil)

	result, err := postService.GetPosts(ctx, page, limit)

	assert.NoError(t, err)
	assert.Equal(t, cachedPosts, result)
	cacheMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestGetPosts_CacheError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		postRepository:        postRepoMock,
		contextService:        contextServiceMock,
	}

	userID := uuid.New()
	page, limit := 1, 10

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("GetPosts", ctx, userID, page, limit).Return(nil, errors.New("cache error"))

	result, err := postService.GetPosts(ctx, page, limit)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cache error")
	assert.Nil(t, result)
	cacheMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestGetPosts_SuccessFromRepositoryAndSetCache(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		postRepository:        postRepoMock,
		contextService:        contextServiceMock,
	}

	userID := uuid.New()
	page, limit := 1, 10
	pagedPosts := &domain.Pagination[*domain.Post]{}

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("GetPosts", ctx, userID, page, limit).Return(nil, nil)
	postRepoMock.On("GetPaginatedPosts", ctx, userID, page, limit).Return(pagedPosts, nil)
	cacheMock.On("SetPost", ctx, userID, mock.Anything, page, limit).Return(nil)

	result, err := postService.GetPosts(ctx, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	cacheMock.AssertExpectations(t)
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestGetPosts_RepositoryError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		postRepository:        postRepoMock,
		contextService:        contextServiceMock,
	}

	userID := uuid.New()
	page, limit := 1, 10

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("GetPosts", ctx, userID, page, limit).Return(nil, nil)
	postRepoMock.On("GetPaginatedPosts", ctx, userID, page, limit).Return(nil, errors.New("repository error"))

	result, err := postService.GetPosts(ctx, page, limit)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository error")
	assert.Nil(t, result)
	cacheMock.AssertExpectations(t)
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestGetPosts_NoPostsFound_ReturnsErrPostNotFound(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		postRepository:        postRepoMock,
		contextService:        contextServiceMock,
	}

	userID := uuid.New()
	page, limit := 1, 10

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("GetPosts", ctx, userID, page, limit).Return(nil, nil)
	postRepoMock.On("GetPaginatedPosts", ctx, userID, page, limit).Return(nil, nil)

	result, err := postService.GetPosts(ctx, page, limit)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrPostNotFound, err)
	assert.Nil(t, result)
	cacheMock.AssertExpectations(t)
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestGetPosts_SetCacheError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		postRepository:        postRepoMock,
		contextService:        contextServiceMock,
	}

	userID := uuid.New()
	page, limit := 1, 10
	pagedPosts := &domain.Pagination[*domain.Post]{}
	pagedPostsResponse := &domain.Pagination[*domain.PostResponse]{
		Limit:      0,
		Page:       0,
		Sort:       "",
		TotalRows:  0,
		TotalPages: 0,
		Rows:       []*domain.PostResponse{},
	}

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("GetPosts", ctx, userID, page, limit).Return(nil, nil)
	postRepoMock.On("GetPaginatedPosts", ctx, userID, page, limit).Return(pagedPosts, nil)
	cacheMock.On("SetPost", ctx, userID, pagedPostsResponse, page, limit).Return(errors.New("cache set error"))

	result, err := postService.GetPosts(ctx, page, limit)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cache set error")
	assert.Nil(t, result)
	cacheMock.AssertExpectations(t)
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestGetPostById_PostNotFound_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	likeRepoMock := new(mocks.LikeRepository)

	postService := &postService{
		postRepository: postRepoMock,
		likeRepository: likeRepoMock,
	}

	postID := uuid.New()
	postRepoMock.On("GetPostById", ctx, postID, true).Return(nil, nil)

	result, err := postService.GetPostById(ctx, postID)

	assert.ErrorIs(t, err, domain.ErrPostNotFound)
	assert.Nil(t, result)
	postRepoMock.AssertExpectations(t)
}

func TestGetPostById_RepositoryError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	likeRepoMock := new(mocks.LikeRepository)

	postService := &postService{
		postRepository: postRepoMock,
		likeRepository: likeRepoMock,
	}

	postID := uuid.New()
	postRepoMock.On("GetPostById", ctx, postID, true).Return(nil, errors.New("repository error"))

	result, err := postService.GetPostById(ctx, postID)

	assert.ErrorContains(t, err, "repository error")
	assert.Nil(t, result)
	postRepoMock.AssertExpectations(t)
}

func TestGetPostById_LikeCheckError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	likeRepoMock := new(mocks.LikeRepository)

	postService := &postService{
		postRepository: postRepoMock,
		likeRepository: likeRepoMock,
	}

	postID := uuid.New()
	userID := uuid.New()
	post := &domain.Post{
		ID:       postID,
		AuthorID: userID,
	}

	postRepoMock.On("GetPostById", ctx, postID, true).Return(post, nil)
	likeRepoMock.On("UserLikedPost", ctx, postID, userID).Return(false, errors.New("like check error"))

	result, err := postService.GetPostById(ctx, postID)

	assert.ErrorContains(t, err, "like check error")
	assert.Nil(t, result)
	postRepoMock.AssertExpectations(t)
	likeRepoMock.AssertExpectations(t)
}

func TestGetPostById_PostFoundAndLiked_ReturnsPostWithLike(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	likeRepoMock := new(mocks.LikeRepository)

	postService := &postService{
		postRepository: postRepoMock,
		likeRepository: likeRepoMock,
	}

	postID := uuid.New()
	userID := uuid.New()
	post := &domain.Post{
		ID:       postID,
		AuthorID: userID,
	}

	postRepoMock.On("GetPostById", ctx, postID, true).Return(post, nil)
	likeRepoMock.On("UserLikedPost", ctx, postID, userID).Return(true, nil)

	result, err := postService.GetPostById(ctx, postID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, postID, result.ID)
	assert.True(t, result.LikesByUser)
	postRepoMock.AssertExpectations(t)
	likeRepoMock.AssertExpectations(t)
}

func TestGetPostById_PostFoundNotLiked_ReturnsPostWithoutLike(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	likeRepoMock := new(mocks.LikeRepository)

	postService := &postService{
		postRepository: postRepoMock,
		likeRepository: likeRepoMock,
	}

	postID := uuid.New()
	userID := uuid.New()
	post := &domain.Post{
		ID:       postID,
		AuthorID: userID,
	}

	postRepoMock.On("GetPostById", ctx, postID, true).Return(post, nil)
	likeRepoMock.On("UserLikedPost", ctx, postID, userID).Return(false, nil)

	result, err := postService.GetPostById(ctx, postID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, postID, result.ID)
	assert.False(t, result.LikesByUser)
	postRepoMock.AssertExpectations(t)
	likeRepoMock.AssertExpectations(t)
}
