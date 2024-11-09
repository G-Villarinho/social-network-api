package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type postService struct {
	di                    *pkg.Di
	postRepository        domain.PostRepository
	contextService        domain.ContextService
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

	return nil, nil
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
	log := slog.With(
		slog.String("service", "post"),
		slog.String("func", "UnlikePost"),
	)

	userID := p.contextService.GetUserID(ctx)

	if err := p.memoryCacheRepository.RemovePostLike(ctx, ID, userID); err != nil {
		return fmt.Errorf("error deleting like from Redis: %w", err)
	}

	go func() {
		message, err := jsoniter.Marshal(domain.LikePayload{UserID: userID, PostID: ID})
		if err != nil {
			log.Error("error to marshal unlike event", slog.String("error", err.Error()))
			return
		}

		if err := p.queueService.Publish(config.QueueUnlikePost, message); err != nil {
			log.Error("error to publish unlike event", slog.String("error", err.Error()))
		}
	}()

	return nil
}

func (p *postService) ProcessLikePost(ctx context.Context, payload domain.LikePayload) error {
	post, err := p.postRepository.GetPostById(ctx, payload.PostID, false)
	if err != nil {
		return fmt.Errorf("post by ID: %w", err)
	}

	if post == nil {
		return domain.ErrPostNotFound
	}

	hasLike, err := p.postRepository.HasUserLikedPost(ctx, payload.PostID, payload.UserID)
	if err != nil {
		return fmt.Errorf("error to check if user has liked post: %w", err)
	}

	if hasLike {
		return domain.ErrPostAlreadyLiked
	}

	if err := p.postRepository.LikePost(ctx, *payload.ToLike()); err != nil {
		return fmt.Errorf("error to like post: %w", err)
	}

	return nil
}

func (p *postService) ProcessUnlikePost(ctx context.Context, payload domain.LikePayload) error {
	post, err := p.postRepository.GetPostById(ctx, payload.PostID, false)
	if err != nil {
		return fmt.Errorf("error to get post by ID: %w", err)
	}

	if post == nil {
		return domain.ErrPostNotFound
	}

	hasLiked, err := p.postRepository.HasUserLikedPost(ctx, payload.UserID, payload.UserID)
	if err != nil {
		return fmt.Errorf("error to check if user has liked post: %w", err)
	}

	if !hasLiked {
		return domain.ErrPostNotLiked
	}

	if err := p.postRepository.UnlikePost(ctx, payload.PostID, payload.UserID); err != nil {
		return fmt.Errorf("unlike post: %w", err)
	}

	return nil
}
