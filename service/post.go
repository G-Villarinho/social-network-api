package service

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/google/uuid"
)

type postService struct {
	di             *pkg.Di
	postRepository domain.PostRepository
}

func NewPostService(di *pkg.Di) (domain.PostService, error) {
	postRepository, err := pkg.Invoke[domain.PostRepository](di)
	if err != nil {
		return nil, err
	}

	return &postService{
		di:             di,
		postRepository: postRepository,
	}, nil
}

func (p *postService) CreatePost(ctx context.Context, payload domain.PostPayload) error {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return domain.ErrSessionNotFound
	}

	post := payload.ToPost(session.UserID)
	if err := p.postRepository.CreatePost(ctx, *post); err != nil {
		return fmt.Errorf("error to create post: %w", err)
	}

	return nil
}

func (p *postService) GetPosts(ctx context.Context) ([]*domain.PostResponse, error) {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return nil, domain.ErrSessionNotFound
	}

	posts, err := p.postRepository.GetPosts(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("error to get posts: %w", err)
	}

	if posts == nil {
		return nil, domain.ErrPostNotFound
	}

	var postsResponse []*domain.PostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToPostResponse())
	}

	return postsResponse, nil
}

func (p *postService) GetPostById(ctx context.Context, ID uuid.UUID) (*domain.PostResponse, error) {
	post, err := p.postRepository.GetPostById(ctx, ID, true)
	if err != nil {
		return nil, fmt.Errorf("error to get post by ID: %w", err)
	}

	if post == nil {
		return nil, domain.ErrPostNotFound
	}

	return post.ToPostResponse(), nil
}

func (p *postService) UpdatePost(ctx context.Context, ID uuid.UUID, payload domain.PostUpdatePayload) error {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return domain.ErrSessionNotFound
	}

	post, err := p.postRepository.GetPostById(ctx, ID, false)
	if err != nil {
		return fmt.Errorf("error to get post by ID: %w", err)
	}

	if post == nil {
		return domain.ErrPostNotFound
	}

	if post.AuthorID != session.UserID {
		return domain.ErrPostNotBelongToUser
	}

	post.Update(payload)
	if err := p.postRepository.UpdatePost(ctx, ID, *post); err != nil {
		return fmt.Errorf("error to update post: %w", err)
	}

	return nil
}

func (p *postService) DeletePost(ctx context.Context, ID uuid.UUID) error {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return domain.ErrSessionNotFound
	}

	post, err := p.postRepository.GetPostById(ctx, ID, false)
	if err != nil {
		return fmt.Errorf("error to get post by ID: %w", err)
	}

	if post == nil {
		return domain.ErrPostNotFound
	}

	if post.AuthorID != session.UserID {
		return domain.ErrPostNotBelongToUser
	}

	if err := p.postRepository.DeletePost(ctx, ID); err != nil {
		return fmt.Errorf("error to delete post: %w", err)
	}

	return nil
}

func (p *postService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.PostResponse, error) {
	posts, err := p.postRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error to get posts by user ID: %w", err)
	}

	if posts == nil {
		return nil, domain.ErrPostNotFound
	}

	var postsResponse []*domain.PostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToPostResponse())
	}

	return postsResponse, nil

}

func (p *postService) LikePost(ctx context.Context, ID uuid.UUID) error {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return domain.ErrSessionNotFound
	}

	post, err := p.postRepository.GetPostById(ctx, ID, false)
	if err != nil {
		return fmt.Errorf("error to get post by ID: %w", err)
	}

	if post == nil {
		return domain.ErrPostNotFound
	}

	hasLiked, err := p.postRepository.HasUserLikedPost(ctx, ID, session.UserID)
	if err != nil {
		return fmt.Errorf("error to check if user has liked post: %w", err)
	}

	if hasLiked {
		return domain.ErrPostAlreadyLiked
	}

	like := domain.Like{
		PostID: ID,
		UserID: session.UserID,
	}

	if err := p.postRepository.LikePost(ctx, like); err != nil {
		return fmt.Errorf("error to like post: %w", err)
	}

	return nil
}

func (p *postService) UnLikePost(ctx context.Context, ID uuid.UUID) error {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return domain.ErrSessionNotFound
	}

	post, err := p.postRepository.GetPostById(ctx, ID, false)
	if err != nil {
		return fmt.Errorf("error to get post by ID: %w", err)
	}

	if post == nil {
		return domain.ErrPostNotFound
	}

	hasLiked, err := p.postRepository.HasUserLikedPost(ctx, ID, session.UserID)
	if err != nil {
		return fmt.Errorf("error to check if user has liked post: %w", err)
	}

	if !hasLiked {
		return domain.ErrPostNotLiked
	}

	if err := p.postRepository.UnLikePost(ctx, ID, session.UserID); err != nil {
		return fmt.Errorf("error to unlike post: %w", err)
	}

	return nil
}
