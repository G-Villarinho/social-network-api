package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/google/uuid"
)

type feedService struct {
	di                    *pkg.Di
	postRepository        domain.PostRepository
	memoryCacheRepository domain.MemoryCacheRepository
	contextService        domain.ContextService
}

func NewFeedService(di *pkg.Di) (domain.FeedService, error) {
	postRepository, err := pkg.Invoke[domain.PostRepository](di)
	if err != nil {
		return nil, err
	}

	memoryCacheRepository, err := pkg.Invoke[domain.MemoryCacheRepository](di)
	if err != nil {
		return nil, err
	}

	contextService, err := pkg.Invoke[domain.ContextService](di)
	if err != nil {
		return nil, err
	}

	return &feedService{
		di:                    di,
		postRepository:        postRepository,
		memoryCacheRepository: memoryCacheRepository,
		contextService:        contextService,
	}, nil
}

func (f *feedService) GenerateFeed(ctx context.Context, page int, limit int) (*domain.Pagination[*domain.PostResponse], error) {
	log := slog.With(
		slog.String("service", "feed"),
		slog.String("func", "GenerateFeed"),
	)

	userID := f.contextService.GetUserID(ctx)

	cachedFeed, err := f.memoryCacheRepository.GetPosts(ctx, userID, page, limit)
	if err == nil && cachedFeed != nil {
		return cachedFeed, nil
	}

	paginatedPosts, err := f.postRepository.GetPaginatedPosts(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("get paginated posts: %w", err)
	}

	postIDMap := make(map[uuid.UUID]*domain.Post)
	for _, post := range paginatedPosts.Rows {
		postIDMap[post.ID] = post
	}

	postIDs := getKeysFromMap(postIDMap)
	likedPostIDs, missingPostIDs, err := f.memoryCacheRepository.GetCachedAndMissingLikes(ctx, userID, postIDs)
	if err != nil {
		return nil, fmt.Errorf("error retrieving likes from Redis: %w", err)
	}

	if len(missingPostIDs) > 0 {
		dbLikedPostIDs, err := f.postRepository.GetLikesByPostIDs(ctx, userID, missingPostIDs)
		if err != nil {
			return nil, fmt.Errorf("error retrieving missing likes from database: %w", err)
		}
		likedPostIDs = append(likedPostIDs, dbLikedPostIDs...)

		if err := f.memoryCacheRepository.SetLikesByPostIDs(ctx, userID, dbLikedPostIDs); err != nil {
			log.Error("error caching likes in Redis", slog.String("error", err.Error()))
		}
	}

	likedPostIDMap := convertToMap(likedPostIDs)
	paginatedResponse, err := domain.Map(paginatedPosts, func(post *domain.Post) *domain.PostResponse {
		likesByUser := likedPostIDMap[post.ID]
		return post.ToPostResponse(likesByUser)
	}), nil

	if err := f.memoryCacheRepository.SetPost(ctx, userID, paginatedResponse, page, limit); err != nil {
		log.Error("error to cache post", slog.String("error", err.Error()))
	}

	return paginatedResponse, nil
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
