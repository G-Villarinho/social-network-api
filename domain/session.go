package domain

//go:generate mockery --name=SessionService --output=../mocks --outpkg=mocks
//go:generate mockery --name=SessionRepository --output=../mocks --outpkg=mocks

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTokenInvalid           = errors.New("invalid token")
	ErrSessionNotFound        = errors.New("token not found")
	ErrorUnexpectedMethod     = errors.New("unexpected signing method")
	ErrTokenNotFoundInContext = errors.New("token not found in context")
	ErrSessionMismatch        = errors.New("session icompatible for user requested")
)

type Session struct {
	UserID    uuid.UUID `json:"id"`
	Token     string    `json:"token"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
}

type SessionService interface {
	CreateSession(ctx context.Context, user User) (string, error)
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	DeleteSession(ctx context.Context, userID uuid.UUID) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session Session) error
	GetSessionByUserID(ctx context.Context, userID uuid.UUID) (*Session, error)
	DeleteSession(ctx context.Context, userId uuid.UUID) error
}
