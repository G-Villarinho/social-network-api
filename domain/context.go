package domain

//go:generate mockery --name=ContextService --output=../mocks --outpkg=mocks

import (
	"context"

	"github.com/google/uuid"
)

type ContextKey string

const (
	SessionKey   ContextKey = "session"
	UserAgentKey ContextKey = "user-agent"
	ClientIPKey  ContextKey = "client-ip"
)

type ContextService interface {
	GetUserID(ctx context.Context) uuid.UUID
	Session(ctx context.Context) (*Session, error)
	GetUserAgent(ctx context.Context) (string, error)
	GetClientIP(ctx context.Context) (string, error)
}
