package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	ErrUserCannotFollowItself = errors.New("user cannot follow itself")
	ErrFollowerNotFound       = errors.New("follower not found")
)

type Follower struct {
	ID         uuid.UUID `gorm:"column:id;type:char(36);primaryKey"`
	UserID     uuid.UUID `gorm:"column:userId;type:char(36);not null"`
	FollowerID uuid.UUID `gorm:"column:followerId;type:char(36);not null"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Follower   User      `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time `gorm:"column:createdAt;not null"`
	UpdatedAt  time.Time `gorm:"column:updatedAt;default:null"`
}

type FollowerHandler interface {
	FollowUser(ctx echo.Context) error
}

type FollowerService interface {
	FollowUser(ctx context.Context, followerId uuid.UUID) error
}

type FollowerRepository interface {
	CreateFollower(ctx context.Context, follower Follower) error
}

func (Follower) TableName() string {
	return "Follower"
}

func (f *Follower) BeforeCreate(tx *gorm.DB) (err error) {
	f.ID = uuid.New()
	f.CreatedAt = time.Now().UTC()
	return
}

func (f *Follower) BeforeUpdate(tx *gorm.DB) (err error) {
	f.UpdatedAt = time.Now().UTC()
	return
}
