package service

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
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
