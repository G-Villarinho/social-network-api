package service

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePost_Success(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	userID := uuid.New()
	payload := domain.PostPayload{
		Content: "This is a test post",
	}

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	postRepoMock.On("CreatePost", ctx, mock.AnythingOfType("domain.Post")).Return(nil)

	err := postService.CreatePost(ctx, payload)

	assert.NoError(t, err)
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestCreatePost_CreatePostError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	userID := uuid.New()
	payload := domain.PostPayload{
		Content: "This is a test post",
	}

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	postRepoMock.On("CreatePost", ctx, mock.AnythingOfType("domain.Post")).Return(errors.New("error creating post"))

	err := postService.CreatePost(ctx, payload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error to create post")
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestGetPosts_CacheHit_ReturnsCachedFeed(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)
	memoryCacheRepoMock := new(mocks.MemoryCacheRepository)

	postService := &postService{
		postRepository:        postRepoMock,
		contextService:        contextServiceMock,
		memoryCacheRepository: memoryCacheRepoMock,
	}

	userID := uuid.New()
	cachedFeed := &domain.Pagination[*domain.PostResponse]{Page: 1, Limit: 10}

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	memoryCacheRepoMock.On("GetPosts", ctx, userID, 1, 10).Return(cachedFeed, nil)

	result, err := postService.GetPosts(ctx, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, cachedFeed, result)
	memoryCacheRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func sortUUIDs(uuids []uuid.UUID) {
	sort.Slice(uuids, func(i, j int) bool {
		return uuids[i].String() < uuids[j].String()
	})
}

func TestGetPosts_CacheMiss_ReturnsRepositoryFeed(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)
	memoryCacheRepoMock := new(mocks.MemoryCacheRepository)

	postService := &postService{
		postRepository:        postRepoMock,
		contextService:        contextServiceMock,
		memoryCacheRepository: memoryCacheRepoMock,
	}

	userID := uuid.New()
	repoFeed := &domain.Pagination[*domain.Post]{
		Page:  1,
		Limit: 3,
		Rows: []*domain.Post{
			{
				ID:        uuid.New(),
				AuthorID:  uuid.New(),
				Author:    domain.User{Username: "test_author_1"},
				Likes:     5,
				Title:     "Test Post 1",
				Content:   "This is the content of test post 1.",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				AuthorID:  uuid.New(),
				Author:    domain.User{Username: "test_author_2"},
				Likes:     8,
				Title:     "Test Post 2",
				Content:   "This is the content of test post 2.",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				AuthorID:  uuid.New(),
				Author:    domain.User{Username: "test_author_3"},
				Likes:     3,
				Title:     "Test Post 3",
				Content:   "This is the content of test post 3.",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		TotalRows:  30,
		TotalPages: 10,
	}

	repoFeedResponse := domain.Map(repoFeed, func(post *domain.Post) *domain.PostResponse {
		return post.ToPostResponse()
	})

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	memoryCacheRepoMock.On("GetPosts", ctx, userID, 1, 10).Return(nil, nil)
	postRepoMock.On("GetPaginatedPosts", ctx, userID, 1, 10).Return(repoFeed, nil)

	missingPostIDs := []uuid.UUID{repoFeed.Rows[0].ID, repoFeed.Rows[1].ID, repoFeed.Rows[2].ID}
	likedPostIDs := []uuid.UUID{}

	memoryCacheRepoMock.On("GetCachedAndMissingLikes", ctx, userID, mock.MatchedBy(func(ids []uuid.UUID) bool {
		sortUUIDs(ids)
		return assert.ElementsMatch(t, ids, missingPostIDs)
	})).Return(likedPostIDs, missingPostIDs, nil)

	postRepoMock.On("GetLikesByPostIDs", ctx, userID, missingPostIDs).Return(likedPostIDs, nil)
	memoryCacheRepoMock.On("SetLikesByPostIDs", ctx, userID, likedPostIDs).Return(nil)

	memoryCacheRepoMock.On("SetPost", ctx, userID, repoFeedResponse, 1, 10).Return(nil)

	result, err := postService.GetPosts(ctx, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, repoFeedResponse, result)
	memoryCacheRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
	postRepoMock.AssertExpectations(t)
}

func TestGetPosts_CacheError_StillFetchesFromRepository(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)
	memoryCacheRepoMock := new(mocks.MemoryCacheRepository)

	postService := &postService{
		postRepository:        postRepoMock,
		contextService:        contextServiceMock,
		memoryCacheRepository: memoryCacheRepoMock,
	}

	userID := uuid.New()
	repoFeed := &domain.Pagination[*domain.PostResponse]{Page: 1, Limit: 10}

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	memoryCacheRepoMock.On("GetPosts", ctx, userID, 1, 10).Return(nil, errors.New("cache error"))
	postRepoMock.On("GetPaginatedPosts", ctx, userID, 1, 10).Return(repoFeed, nil)

	result, err := postService.GetPosts(ctx, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, repoFeed, result)
	memoryCacheRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
	postRepoMock.AssertExpectations(t)
}

func TestGetPosts_RepositoryError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)
	memoryCacheRepoMock := new(mocks.MemoryCacheRepository)

	postService := &postService{
		postRepository:        postRepoMock,
		contextService:        contextServiceMock,
		memoryCacheRepository: memoryCacheRepoMock,
	}

	userID := uuid.New()

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	memoryCacheRepoMock.On("GetPosts", ctx, userID, 1, 10).Return(nil, nil)
	postRepoMock.On("GetPaginatedPosts", ctx, userID, 1, 10).Return(nil, errors.New("repository error"))

	result, err := postService.GetPosts(ctx, 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "repository error")
	memoryCacheRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
	postRepoMock.AssertExpectations(t)
}

func TestGetPosts_SetCacheError_ReturnsRepositoryFeed(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)
	memoryCacheRepoMock := new(mocks.MemoryCacheRepository)

	postService := &postService{
		postRepository:        postRepoMock,
		contextService:        contextServiceMock,
		memoryCacheRepository: memoryCacheRepoMock,
	}

	userID := uuid.New()
	repoFeed := &domain.Pagination[*domain.PostResponse]{Page: 1, Limit: 10}

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	memoryCacheRepoMock.On("GetPosts", ctx, userID, 1, 10).Return(nil, nil)
	postRepoMock.On("GetPaginatedPosts", ctx, userID, 1, 10).Return(repoFeed, nil)
	memoryCacheRepoMock.On("SetPost", ctx, userID, repoFeed, 1, 10).Return(errors.New("cache set error"))

	result, err := postService.GetPosts(ctx, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, repoFeed, result)
	memoryCacheRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
	postRepoMock.AssertExpectations(t)
}
