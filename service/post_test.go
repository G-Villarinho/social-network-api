package service

import (
	"context"
	"errors"
	"testing"

	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPosts_WhenSuccessFromCache_ShouldReturnPosts(t *testing.T) {
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

func TestGetPosts_WhenCacheError_ShouldReturnError(t *testing.T) {
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

func TestGetPosts_WhenSuccessFromRepositoryAndSetCache_ShouldReturnPosts(t *testing.T) {
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

func TestGetPosts_WhenRepositoryError_ShouldReturnError(t *testing.T) {
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

func TestGetPosts_WhenNoPostsFound_ShouldReturnErrPostNotFound(t *testing.T) {
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

func TestGetPosts_WhenSetCacheError_ShouldReturnError(t *testing.T) {
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

func TestUpdatePost_PostNotFound_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	postID := uuid.New()
	payload := domain.PostUpdatePayload{}

	postRepoMock.On("GetPostById", ctx, postID, false).Return(nil, nil)

	err := postService.UpdatePost(ctx, postID, payload)

	assert.ErrorIs(t, err, domain.ErrPostNotFound)
	postRepoMock.AssertExpectations(t)
}

func TestUpdatePost_GetPostByIdRepositoryError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	postID := uuid.New()
	payload := domain.PostUpdatePayload{}

	postRepoMock.On("GetPostById", ctx, postID, false).Return(nil, errors.New("repository error"))

	err := postService.UpdatePost(ctx, postID, payload)

	assert.ErrorContains(t, err, "repository error")
	postRepoMock.AssertExpectations(t)
}

func TestUpdatePost_PostNotBelongToUser_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()
	payload := domain.PostUpdatePayload{}
	post := &domain.Post{AuthorID: otherUserID}

	postRepoMock.On("GetPostById", ctx, postID, false).Return(post, nil)
	contextServiceMock.On("GetUserID", ctx).Return(userID)

	err := postService.UpdatePost(ctx, postID, payload)

	assert.ErrorIs(t, err, domain.ErrPostNotBelongToUser)
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestUpdatePost_UpdatePostRepositoryError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()
	payload := domain.PostUpdatePayload{}
	post := &domain.Post{AuthorID: userID}

	postRepoMock.On("GetPostById", ctx, postID, false).Return(post, nil)
	contextServiceMock.On("GetUserID", ctx).Return(userID)
	postRepoMock.On("UpdatePost", ctx, postID, *post).Return(errors.New("repository update error"))

	err := postService.UpdatePost(ctx, postID, payload)

	assert.ErrorContains(t, err, "repository update error")
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestUpdatePost_Success(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()
	payload := domain.PostUpdatePayload{}
	post := &domain.Post{AuthorID: userID}

	postRepoMock.On("GetPostById", ctx, postID, false).Return(post, nil)
	contextServiceMock.On("GetUserID", ctx).Return(userID)
	postRepoMock.On("UpdatePost", ctx, postID, *post).Return(nil)

	err := postService.UpdatePost(ctx, postID, payload)

	assert.NoError(t, err)
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestDeletePost_PostNotFound_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	postID := uuid.New()

	postRepoMock.On("GetPostById", ctx, postID, false).Return(nil, nil)

	err := postService.DeletePost(ctx, postID)

	assert.ErrorIs(t, err, domain.ErrPostNotFound)
	postRepoMock.AssertExpectations(t)
}

func TestDeletePost_GetPostByIdRepositoryError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	postID := uuid.New()

	postRepoMock.On("GetPostById", ctx, postID, false).Return(nil, errors.New("repository error"))

	err := postService.DeletePost(ctx, postID)

	assert.ErrorContains(t, err, "repository error")
	postRepoMock.AssertExpectations(t)
}

func TestDeletePost_PostNotBelongToUser_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()
	post := &domain.Post{AuthorID: otherUserID}

	postRepoMock.On("GetPostById", ctx, postID, false).Return(post, nil)
	contextServiceMock.On("GetUserID", ctx).Return(userID)

	err := postService.DeletePost(ctx, postID)

	assert.ErrorIs(t, err, domain.ErrPostNotBelongToUser)
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestDeletePost_DeletePostRepositoryError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()
	post := &domain.Post{AuthorID: userID}

	postRepoMock.On("GetPostById", ctx, postID, false).Return(post, nil)
	contextServiceMock.On("GetUserID", ctx).Return(userID)
	postRepoMock.On("DeletePost", ctx, postID).Return(errors.New("delete error"))

	err := postService.DeletePost(ctx, postID)

	assert.ErrorContains(t, err, "delete error")
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestDeletePost_Success(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	contextServiceMock := new(mocks.ContextService)

	postService := &postService{
		postRepository: postRepoMock,
		contextService: contextServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()
	post := &domain.Post{AuthorID: userID}

	postRepoMock.On("GetPostById", ctx, postID, false).Return(post, nil)
	contextServiceMock.On("GetUserID", ctx).Return(userID)
	postRepoMock.On("DeletePost", ctx, postID).Return(nil)

	err := postService.DeletePost(ctx, postID)

	assert.NoError(t, err)
	postRepoMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestGetByUserID_SessionNotFound_ReturnsError(t *testing.T) {
	ctx := context.Background()
	postRepoMock := new(mocks.PostRepository)
	likeRepoMock := new(mocks.LikeRepository)

	postService := &postService{
		postRepository: postRepoMock,
		likeRepository: likeRepoMock,
	}

	userID := uuid.New()

	postsResponse, err := postService.GetByUserID(ctx, userID)

	assert.ErrorIs(t, err, domain.ErrSessionNotFound)
	assert.Nil(t, postsResponse)
	postRepoMock.AssertNotCalled(t, "GetByUserID", mock.Anything, mock.Anything)
	likeRepoMock.AssertNotCalled(t, "GetLikedPostIDs", mock.Anything, mock.Anything)
}

func TestGetByUserID_GetByUserIDRepositoryError_ReturnsError(t *testing.T) {
	session := &domain.Session{UserID: uuid.New()}
	ctx := context.WithValue(context.Background(), domain.SessionKey, session)
	postRepoMock := new(mocks.PostRepository)
	likeRepoMock := new(mocks.LikeRepository)

	postService := &postService{
		postRepository: postRepoMock,
		likeRepository: likeRepoMock,
	}

	userID := uuid.New()

	postRepoMock.On("GetByUserID", ctx, userID).Return(nil, errors.New("repository error"))

	postsResponse, err := postService.GetByUserID(ctx, userID)

	assert.ErrorContains(t, err, "repository error")
	assert.Nil(t, postsResponse)
	postRepoMock.AssertExpectations(t)
	likeRepoMock.AssertNotCalled(t, "GetLikedPostIDs", mock.Anything, mock.Anything)
}

func TestGetByUserID_PostsNotFound_ReturnsError(t *testing.T) {
	session := &domain.Session{UserID: uuid.New()}
	ctx := context.WithValue(context.Background(), domain.SessionKey, session)
	postRepoMock := new(mocks.PostRepository)
	likeRepoMock := new(mocks.LikeRepository)

	postService := &postService{
		postRepository: postRepoMock,
		likeRepository: likeRepoMock,
	}

	userID := uuid.New()

	postRepoMock.On("GetByUserID", ctx, userID).Return(nil, nil)

	postsResponse, err := postService.GetByUserID(ctx, userID)

	assert.ErrorIs(t, err, domain.ErrPostNotFound)
	assert.Nil(t, postsResponse)
	postRepoMock.AssertExpectations(t)
	likeRepoMock.AssertNotCalled(t, "GetLikedPostIDs", mock.Anything, mock.Anything)
}

func TestGetByUserID_GetLikedPostIDsError_ReturnsError(t *testing.T) {
	session := &domain.Session{UserID: uuid.New()}
	ctx := context.WithValue(context.Background(), domain.SessionKey, session)
	postRepoMock := new(mocks.PostRepository)
	likeRepoMock := new(mocks.LikeRepository)

	postService := &postService{
		postRepository: postRepoMock,
		likeRepository: likeRepoMock,
	}

	userID := uuid.New()
	posts := []*domain.Post{{ID: uuid.New()}}

	postRepoMock.On("GetByUserID", ctx, userID).Return(posts, nil)
	likeRepoMock.On("GetLikedPostIDs", ctx, session.UserID).Return(nil, errors.New("like repository error"))

	postsResponse, err := postService.GetByUserID(ctx, userID)

	assert.ErrorContains(t, err, "like repository error")
	assert.Nil(t, postsResponse)
	postRepoMock.AssertExpectations(t)
	likeRepoMock.AssertExpectations(t)
}

func TestGetByUserID_Success_ReturnsPostsWithLikes(t *testing.T) {
	session := &domain.Session{UserID: uuid.New()}
	ctx := context.WithValue(context.Background(), domain.SessionKey, session)
	postRepoMock := new(mocks.PostRepository)
	likeRepoMock := new(mocks.LikeRepository)

	postService := &postService{
		postRepository: postRepoMock,
		likeRepository: likeRepoMock,
	}

	userID := uuid.New()
	postID := uuid.New()
	posts := []*domain.Post{{ID: postID}}
	likedPostIDs := map[uuid.UUID]bool{postID: true}

	postRepoMock.On("GetByUserID", ctx, userID).Return(posts, nil)
	likeRepoMock.On("GetLikedPostIDs", ctx, session.UserID).Return(likedPostIDs, nil)

	postsResponse, err := postService.GetByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, postsResponse)
	assert.Len(t, postsResponse, 1)
	assert.True(t, postsResponse[0].LikesByUser)
	postRepoMock.AssertExpectations(t)
	likeRepoMock.AssertExpectations(t)
}

func TestGetByUserID_Success_ReturnsPostsWithoutLikes(t *testing.T) {
	session := &domain.Session{UserID: uuid.New()}
	ctx := context.WithValue(context.Background(), domain.SessionKey, session)
	postRepoMock := new(mocks.PostRepository)
	likeRepoMock := new(mocks.LikeRepository)

	postService := &postService{
		postRepository: postRepoMock,
		likeRepository: likeRepoMock,
	}

	userID := uuid.New()
	postID := uuid.New()
	posts := []*domain.Post{{ID: postID}}
	likedPostIDs := map[uuid.UUID]bool{}

	postRepoMock.On("GetByUserID", ctx, userID).Return(posts, nil)
	likeRepoMock.On("GetLikedPostIDs", ctx, session.UserID).Return(likedPostIDs, nil)

	postsResponse, err := postService.GetByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, postsResponse)
	assert.Len(t, postsResponse, 1)
	assert.False(t, postsResponse[0].LikesByUser)
	postRepoMock.AssertExpectations(t)
	likeRepoMock.AssertExpectations(t)
}

func TestLikePost_Success(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	contextServiceMock := new(mocks.ContextService)
	queueServiceMock := new(mocks.QueueService)

	done := make(chan bool)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		contextService:        contextServiceMock,
		queueService:          queueServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("SetPostLike", ctx, postID, userID).Return(nil)
	queueServiceMock.On("Publish", config.QueueLikePost, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done <- true
	})

	err := postService.LikePost(ctx, postID)

	<-done

	assert.NoError(t, err)
	cacheMock.AssertExpectations(t)
	queueServiceMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestLikePost_CacheError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	contextServiceMock := new(mocks.ContextService)
	queueServiceMock := new(mocks.QueueService)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		contextService:        contextServiceMock,
		queueService:          queueServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("SetPostLike", ctx, postID, userID).Return(errors.New("cache error"))

	err := postService.LikePost(ctx, postID)

	assert.ErrorContains(t, err, "cache error")
	cacheMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestLikePost_QueuePublishError_LogError(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	contextServiceMock := new(mocks.ContextService)
	queueServiceMock := new(mocks.QueueService)
	done := make(chan bool)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		contextService:        contextServiceMock,
		queueService:          queueServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("SetPostLike", ctx, postID, userID).Return(nil)
	queueServiceMock.On("Publish", config.QueueLikePost, mock.Anything).Return(errors.New("publish error")).Run(func(args mock.Arguments) {
		done <- true
	})

	err := postService.LikePost(ctx, postID)

	<-done

	assert.NoError(t, err)
	cacheMock.AssertExpectations(t)
	queueServiceMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestUnlikePost_Success(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	contextServiceMock := new(mocks.ContextService)
	queueServiceMock := new(mocks.QueueService)

	done := make(chan bool)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		contextService:        contextServiceMock,
		queueService:          queueServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("RemovePostLike", ctx, postID, userID).Return(nil)
	queueServiceMock.On("Publish", config.QueueUnlikePost, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done <- true
	})

	err := postService.UnlikePost(ctx, postID)

	<-done

	assert.NoError(t, err)
	cacheMock.AssertExpectations(t)
	queueServiceMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestUnlikePost_CacheError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	contextServiceMock := new(mocks.ContextService)
	queueServiceMock := new(mocks.QueueService)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		contextService:        contextServiceMock,
		queueService:          queueServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("RemovePostLike", ctx, postID, userID).Return(errors.New("cache error"))

	err := postService.UnlikePost(ctx, postID)

	assert.ErrorContains(t, err, "cache error")
	cacheMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}

func TestUnlikePost_QueuePublishError_LogError(t *testing.T) {
	ctx := context.Background()
	cacheMock := new(mocks.MemoryCacheRepository)
	contextServiceMock := new(mocks.ContextService)
	queueServiceMock := new(mocks.QueueService)

	done := make(chan bool)

	postService := &postService{
		memoryCacheRepository: cacheMock,
		contextService:        contextServiceMock,
		queueService:          queueServiceMock,
	}

	postID := uuid.New()
	userID := uuid.New()

	contextServiceMock.On("GetUserID", ctx).Return(userID)
	cacheMock.On("RemovePostLike", ctx, postID, userID).Return(nil)
	queueServiceMock.On("Publish", config.QueueUnlikePost, mock.Anything).Return(errors.New("publish error")).Run(func(args mock.Arguments) {
		done <- true
	})

	err := postService.UnlikePost(ctx, postID)

	<-done

	assert.NoError(t, err)
	cacheMock.AssertExpectations(t)
	queueServiceMock.AssertExpectations(t)
	contextServiceMock.AssertExpectations(t)
}
