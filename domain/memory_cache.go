package domain

//go:generate mockery --name=MemoryCacheRepository --dir=. --output=../mocks/ --outpkg=mocks

import (
	"context"

	"github.com/google/uuid"
)

type MemoryCacheRepository interface {
	SetPostLike(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
	RemovePostLike(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
	SetPost(ctx context.Context, userID uuid.UUID, posts *Pagination[*PostResponse], page, limit int) error
	GetPosts(ctx context.Context, userID uuid.UUID, page, limit int) (*Pagination[*PostResponse], error)
	GetCachedAndMissingLikes(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) ([]uuid.UUID, []uuid.UUID, error)
	SetLikesByPostIDs(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) error
}
