package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/G-Villarinho/social-network/client"
	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type postService struct {
	di                    *pkg.Di
	redisClient           *redis.Client
	postRepository        domain.PostRepository
	contextService        domain.ContextService
	rabbitMQClient        client.RabbitMQClient
	memoryCacheRepository domain.MemoryCacheRepository
	queueService          domain.QueueService
}

func NewPostService(di *pkg.Di) (domain.PostService, error) {
	postRepository, err := pkg.Invoke[domain.PostRepository](di)
	if err != nil {
		return nil, err
	}

	contextService, err := pkg.Invoke[domain.ContextService](di)
	if err != nil {
		return nil, err
	}

	redisClient, err := pkg.Invoke[*redis.Client](di)
	if err != nil {
		return nil, err
	}

	rabbitMQClient, err := pkg.Invoke[client.RabbitMQClient](di)
	if err != nil {
		return nil, err
	}

	memoryCacheRepository, err := pkg.Invoke[domain.MemoryCacheRepository](di)
	if err != nil {
		return nil, err
	}

	queueService, err := pkg.Invoke[domain.QueueService](di)
	if err != nil {
		return nil, err
	}

	return &postService{
		di:                    di,
		postRepository:        postRepository,
		contextService:        contextService,
		redisClient:           redisClient,
		rabbitMQClient:        rabbitMQClient,
		memoryCacheRepository: memoryCacheRepository,
		queueService:          queueService,
	}, nil
}

func (p *postService) CreatePost(ctx context.Context, payload domain.PostPayload) error {
	post := payload.ToPost(p.contextService.GetUserID(ctx))
	if err := p.postRepository.CreatePost(ctx, *post); err != nil {
		return fmt.Errorf("error to create post: %w", err)
	}

	return nil
}

func (p *postService) GetPosts(ctx context.Context, page int, limit int) (*domain.Pagination[*domain.PostResponse], error) {
	cachedFeed, err := p.getCachedPosts(ctx, page, limit)
	if err == nil && cachedFeed != nil {
		return cachedFeed, nil
	}

	paginatedPosts, err := p.buildPaginatedResponse(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	if err := p.cachePosts(ctx, paginatedPosts, page, limit); err != nil {
		slog.Error("error to cache post", slog.String("error", err.Error()))
	}

	return paginatedPosts, nil
}

func (p *postService) GetPostById(ctx context.Context, ID uuid.UUID) (*domain.PostResponse, error) {
	post, err := p.postRepository.GetPostById(ctx, ID, true)
	if err != nil {
		return nil, fmt.Errorf("error to get post by ID: %w", err)
	}

	if post == nil {
		return nil, domain.ErrPostNotFound
	}

	hasLiked, err := p.postRepository.HasUserLikedPost(ctx, ID, post.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("error to check if user has liked post: %w", err)
	}

	return post.ToPostResponse(hasLiked), nil
}

func (p *postService) UpdatePost(ctx context.Context, ID uuid.UUID, payload domain.PostUpdatePayload) error {
	post, err := p.postRepository.GetPostById(ctx, ID, false)
	if err != nil {
		return fmt.Errorf("error to get post by ID: %w", err)
	}

	if post == nil {
		return domain.ErrPostNotFound
	}

	if post.AuthorID != p.contextService.GetUserID(ctx) {
		return domain.ErrPostNotBelongToUser
	}

	post.Update(payload)
	if err := p.postRepository.UpdatePost(ctx, ID, *post); err != nil {
		return fmt.Errorf("error to update post: %w", err)
	}

	return nil
}

func (p *postService) DeletePost(ctx context.Context, ID uuid.UUID) error {
	post, err := p.postRepository.GetPostById(ctx, ID, false)
	if err != nil {
		return fmt.Errorf("error to get post by ID: %w", err)
	}

	if post == nil {
		return domain.ErrPostNotFound
	}

	if post.AuthorID != p.contextService.GetUserID(ctx) {
		return domain.ErrPostNotBelongToUser
	}

	if err := p.postRepository.DeletePost(ctx, ID); err != nil {
		return fmt.Errorf("error to delete post: %w", err)
	}

	return nil
}

func (p *postService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.PostResponse, error) {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return nil, domain.ErrSessionNotFound
	}

	posts, err := p.postRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error to get posts by user ID: %w", err)
	}

	if posts == nil {
		return nil, domain.ErrPostNotFound
	}

	likedPostIDs, err := p.postRepository.GetLikedPostIDs(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	var postsResponse []*domain.PostResponse
	for _, post := range posts {
		likesByUser := false
		if _, exists := likedPostIDs[post.ID]; exists {
			likesByUser = true
		}

		postsResponse = append(postsResponse, post.ToPostResponse(likesByUser))
	}

	return postsResponse, nil
}

func (p *postService) LikePost(ctx context.Context, ID uuid.UUID) error {
	log := slog.With(
		slog.String("service", "post"),
		slog.String("func", "LikePost"),
	)

	userID := p.contextService.GetUserID(ctx)

	if err := p.memoryCacheRepository.SetPostLike(ctx, ID, userID); err != nil {
		return fmt.Errorf("error caching like in Redis: %w", err)
	}

	go func() {
		message, err := jsoniter.Marshal(domain.LikePayload{UserID: userID, PostID: ID})
		if err != nil {
			log.Error("error to marshal like event", slog.String("error", err.Error()))
			return
		}

		if err := p.queueService.Publish(config.QueueLikePost, message); err != nil {
			log.Error("error to publish like event", slog.String("error", err.Error()))
		}
	}()

	return nil
}

func (p *postService) UnlikePost(ctx context.Context, ID uuid.UUID) error {
	post, err := p.postRepository.GetPostById(ctx, ID, false)
	if err != nil {
		return fmt.Errorf("error to get post by ID: %w", err)
	}

	if post == nil {
		return domain.ErrPostNotFound
	}

	hasLiked, err := p.postRepository.HasUserLikedPost(ctx, ID, p.contextService.GetUserID(ctx))
	if err != nil {
		return fmt.Errorf("error to check if user has liked post: %w", err)
	}

	if !hasLiked {
		return domain.ErrPostNotLiked
	}

	if err := p.postRepository.UnlikePost(ctx, ID, p.contextService.GetUserID(ctx)); err != nil {
		return fmt.Errorf("error to unlike post: %w", err)
	}

	return nil
}

func (p *postService) buildPaginatedResponse(ctx context.Context, page, limit int) (*domain.Pagination[*domain.PostResponse], error) {
	userID := p.contextService.GetUserID(ctx)

	paginatedPosts, err := p.postRepository.GetPaginatedPosts(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("error to get paginated posts: %w", err)
	}

	postIDMap := make(map[uuid.UUID]*domain.Post)
	for _, post := range paginatedPosts.Rows {
		postIDMap[post.ID] = post
	}

	likedPostIDs, err := p.postRepository.GetLikesByPostIDs(ctx, userID, getKeysFromMap(postIDMap))
	if err != nil {
		return nil, err
	}

	likedPostIDMap := convertToMap(likedPostIDs)

	paginatedResponse := &domain.Pagination[*domain.PostResponse]{
		Limit:      paginatedPosts.Limit,
		Page:       paginatedPosts.Page,
		TotalRows:  paginatedPosts.TotalRows,
		TotalPages: paginatedPosts.TotalPages,
		Rows:       make([]*domain.PostResponse, 0, len(paginatedPosts.Rows)),
	}

	for _, post := range paginatedPosts.Rows {
		likesByUser := likedPostIDMap[post.ID]
		paginatedResponse.Rows = append(paginatedResponse.Rows, post.ToPostResponse(likesByUser))
	}

	return paginatedResponse, nil
}

func (p *postService) getCachedPosts(ctx context.Context, page, limit int) (*domain.Pagination[*domain.PostResponse], error) {
	var cachedFeed domain.Pagination[*domain.PostResponse]

	data, err := p.redisClient.Get(ctx, getKeyCachePost(p.contextService.GetUserID(ctx), page, limit)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}
	if err := jsoniter.UnmarshalFromString(data, &cachedFeed); err != nil {
		return nil, err
	}
	return &cachedFeed, nil
}

func (p *postService) cachePosts(ctx context.Context, feed *domain.Pagination[*domain.PostResponse], page, limit int) error {
	data, err := jsoniter.MarshalToString(feed)
	if err != nil {
		return err
	}

	return p.redisClient.Set(ctx, getKeyCachePost(p.contextService.GetUserID(ctx), page, limit), data, 5*time.Minute).Err()
}

func getKeysFromMap(m map[uuid.UUID]*domain.Post) []uuid.UUID {
	keys := make([]uuid.UUID, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func convertToMap(ids []uuid.UUID) map[uuid.UUID]bool {
	m := make(map[uuid.UUID]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m
}

func getKeyCachePost(userID uuid.UUID, page, limit int) string {
	return fmt.Sprintf("user:%s:feed:page:%d:limit:%d", userID, page, limit)
}
