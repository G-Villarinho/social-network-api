package service

import (
	"context"
	"fmt"
	"log/slog"

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

func (p *postService) GetPosts(ctx context.Context, page int, limit int) (*domain.Pagination[*domain.PostResponse], error) {
	log := slog.With(
		slog.String("service", "post"),
		slog.String("func", "GetPosts"),
	)

	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return nil, domain.ErrSessionNotFound
	}

	cachedFeed, err := p.postRepository.GetCachedPosts(ctx, getKeyCachePost(session.UserID, page, limit))
	if err == nil && cachedFeed != nil {
		return cachedFeed, nil
	}

	paginatedPosts, err := p.postRepository.GetPaginatedPosts(ctx, session.UserID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("error to get paginated feed: %w", err)
	}

	postIDMap := make(map[uuid.UUID]*domain.Post)
	for _, post := range paginatedPosts.Rows {
		postIDMap[post.ID] = post
	}

	likedPostIDs, err := p.postRepository.GetLikesByPostIDs(ctx, session.UserID, getKeysFromMap(postIDMap))
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

	if err := p.postRepository.CachePost(ctx, getKeyCachePost(session.UserID, page, limit), paginatedResponse); err != nil {
		log.Error("error to cache post", slog.String("error", err.Error()))
	}

	return paginatedResponse, nil
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

// TODO: Implement LikePost and UnlikePost methods asynchonously
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

func (p *postService) UnlikePost(ctx context.Context, ID uuid.UUID) error {
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

	if err := p.postRepository.UnlikePost(ctx, ID, session.UserID); err != nil {
		return fmt.Errorf("error to unlike post: %w", err)
	}

	return nil
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
