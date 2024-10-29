package service

import (
	"context"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/google/uuid"
)

type contextService struct {
	di *pkg.Di
}

func NewContextService(di *pkg.Di) (domain.ContextService, error) {
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
