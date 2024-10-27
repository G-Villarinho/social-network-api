package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

type sessionRepository struct {
	di          *pkg.Di
	db          *gorm.DB
	redisClient *redis.Client
}

func NewSessionRepository(di *pkg.Di) (domain.SessionRepository, error) {
	db, err := pkg.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	redisClient, err := pkg.Invoke[*redis.Client](di)
	if err != nil {
		return nil, err
	}

	return &sessionRepository{
		di:          di,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (s *sessionRepository) CreateSession(ctx context.Context, session domain.Session) error {
	sessionJSON, err := jsoniter.Marshal(session)
	if err != nil {
		return err
	}

	if err := s.redisClient.Set(ctx, s.getSessionKey(session.UserID.String()), sessionJSON, time.Duration(config.Env.SessionExp)*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

func (s *sessionRepository) GetSessionByUserID(ctx context.Context, userID uuid.UUID) (*domain.Session, error) {
	sessionJSON, err := s.redisClient.Get(ctx, s.getSessionKey(userID.String())).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var session domain.Session
	if err := jsoniter.UnmarshalFromString(sessionJSON, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *sessionRepository) getSessionKey(userID string) string {
	return fmt.Sprintf("session_%s", userID)
}
