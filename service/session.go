package service

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/golang-jwt/jwt"
	jsoniter "github.com/json-iterator/go"
)

type sessionService struct {
	di                *pkg.Di
	sessionRepository domain.SessionRepository
}

func NewSessionService(di *pkg.Di) (domain.SessionService, error) {
	sessionRepository, err := pkg.Invoke[domain.SessionRepository](di)
	if err != nil {
		return nil, err
	}

	return &sessionService{
		di:                di,
		sessionRepository: sessionRepository,
	}, nil
}

func (s *sessionService) CreateSession(ctx context.Context, user domain.User) (string, error) {
	token, err := s.createToken(user)
	if err != nil {
		return "", fmt.Errorf("error to create token for user ID %s: %w", user.ID, err)
	}

	session := &domain.Session{
		UserID:    user.ID,
		Token:     token,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    user.Avatar,
	}

	if err := s.sessionRepository.CreateSession(ctx, *session); err != nil {
		return "", fmt.Errorf("error to create session for user ID %s: %w", user.ID, err)
	}

	return token, nil

}

func (s *sessionService) GetSessionByToken(ctx context.Context, token string) (*domain.Session, error) {
	sessionFromToken, err := s.extractSessionFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("error to extract session from token: %w", err)
	}

	session, err := s.sessionRepository.GetSessionByUserID(ctx, sessionFromToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("error to get session for user ID %s: %w", sessionFromToken.UserID, err)
	}

	if session == nil {
		return nil, domain.ErrSessionNotFound
	}

	if token != session.Token {
		return nil, domain.ErrSessionMismatch
	}

	return session, nil
}

func (s *sessionService) createToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"id":        user.ID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
		"avatar":    user.Avatar,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	tokenString, err := token.SignedString(config.Env.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("error to sign token for user ID %s: %w", user.ID, err)
	}

	return tokenString, nil
}

func (s *sessionService) extractSessionFromToken(tokenString string) (*domain.Session, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, domain.ErrorUnexpectedMethod
		}
		return config.Env.PublicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error to parse token: %w", err)
	}

	if !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrTokenInvalid
	}

	sessionJSON, err := jsoniter.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("error to marshal claims into JSON: %w", err)
	}

	var session domain.Session
	if err := jsoniter.Unmarshal(sessionJSON, &session); err != nil {
		return nil, fmt.Errorf("error to unmarshal session from JSON: %w", err)
	}

	return &session, nil
}
