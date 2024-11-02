package domain

import (
	"context"

	"github.com/google/uuid"
)

type MemoryCacheRepository interface {
	SetPostLike(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
	RemovePostLike(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
}
