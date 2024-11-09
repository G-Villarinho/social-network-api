package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

//go:generate mockery --name=FollowerHandler --output=../mocks --outpkg=mocks
//go:generate mockery --name=FollowerService --output=../mocks --outpkg=mocks
//go:generate mockery --name=FollowerRepository --output=../mocks --outpkg=mocks

var (
	ErrUserCannotFollowItself   = errors.New("user cannot follow itself")
	ErrFollowerNotFound         = errors.New("follower not found")
	ErrorFollowerAlreadyExist   = errors.New("follower already exist")
	ErrUserCannotUnfollowItself = errors.New("user cannot unfollow itself")
	ErrFollowerAlreadyExists    = errors.New("follower already exists")
	ErrFollowingNotFound        = errors.New("following not found")
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

type FollowerResponse struct {
	ID        uuid.UUID             `json:"id"`
	User      *UserFollowerResponse `json:"user"`
	CreatedAt time.Time             `json:"createdAt"`
}

type FollowerHandler interface {
	FollowUser(ctx echo.Context) error
	UnfollowUser(ctx echo.Context) error
	GetFollowers(ctx echo.Context) error
	GetFollowings(ctx echo.Context) error
}

type FollowerService interface {
	FollowUser(ctx context.Context, userID uuid.UUID) error
	UnfollowUser(ctx context.Context, userID uuid.UUID) error
	GetFollowers(ctx context.Context) ([]*FollowerResponse, error)
	GetFollowings(ctx context.Context) ([]*FollowerResponse, error)
}

type FollowerRepository interface {
	CreateFollower(ctx context.Context, follower Follower) error
	DeleteFollower(ctx context.Context, followerId uuid.UUID) error
	GetFollower(ctx context.Context, userID uuid.UUID, followerId uuid.UUID) (*Follower, error)
	GetFollowers(ctx context.Context, userID uuid.UUID) ([]*Follower, error)
	GetFollowings(ctx context.Context, userID uuid.UUID) ([]*Follower, error)
}

func (f *Follower) ToFollowerResponse() *FollowerResponse {
	response := &FollowerResponse{
		ID:        f.ID,
		CreatedAt: f.CreatedAt,
	}

	if f.Follower != (User{}) {
		response.User = f.Follower.ToUserFollowerResponse()
	} else if f.User != (User{}) {
		response.User = f.User.ToUserFollowerResponse()
	}

	return response
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
