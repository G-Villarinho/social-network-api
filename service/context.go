package service

import (
	"context"
	"errors"

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

func (c *contextService) GetUserAgent(ctx context.Context) (string, error) {
	userAgent, ok := ctx.Value(domain.UserAgentKey).(string)
	if !ok {
		return "", errors.New("user-agent not found in context")
	}
	return userAgent, nil
}

func (c *contextService) GetClientIP(ctx context.Context) (string, error) {
	clientIP, ok := ctx.Value(domain.ClientIPKey).(string)
	if !ok {
		return "", errors.New("client-ip not found in context")
	}
	return clientIP, nil
}
