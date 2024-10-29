package domain

import (
	"context"

	"github.com/google/uuid"
)

type ContextKey string

const (
	SessionKey ContextKey = "session"
)

type ContextService interface {
	GetUserID(ctx context.Context) uuid.UUID
	Session(ctx context.Context) (*Session, error)
}
