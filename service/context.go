package service

import (
	"context"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/google/uuid"
)

type contextService struct {
	di *internal.Di
}

func NewContextService(di *internal.Di) (domain.ContextService, error) {
	return &contextService{di: di}, nil
}

func (c *contextService) Session(ctx context.Context) (*domain.Session, error) {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return nil, domain.ErrSessionNotFound
	}
	return session, nil
}

func (c *contextService) GetUserID(ctx context.Context) uuid.UUID {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		panic("session not found")
	}

	return session.UserID
}
