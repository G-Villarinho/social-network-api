package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Like struct {
	ID        uuid.UUID `gorm:"column:id;type:char(36);primaryKey"`
	UserID    uuid.UUID `gorm:"column:userID;type:char(36);not null"`
	PostID    uuid.UUID `gorm:"column:postID;type:char(36);not null"`
	Post      Post      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"column:createdAt;not null"`
	UpdatedAt time.Time `gorm:"column:updatedAt;default:null"`
}

type LikePayload struct {
	UserID uuid.UUID `json:"userId"`
	PostID uuid.UUID `json:"postId"`
}

type LikeService interface {
	CreateLike(ctx context.Context, payload LikePayload) error
}

type LikeRepository interface {
	CreateLike(ctx context.Context, like Like) error
	UserLikedPost(ctx context.Context, ID uuid.UUID, userID uuid.UUID) (bool, error)
}

func (l *LikePayload) ToLike() *Like {
	return &Like{
		UserID: l.UserID,
		PostID: l.PostID,
	}
}

func (Like) TableName() string {
	return "Like"
}

func (l *Like) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New()
	l.CreatedAt = time.Now().UTC()
	return
}

func (l *Like) BeforeUpdate(tx *gorm.DB) (err error) {
	l.UpdatedAt = time.Now().UTC()
	return
}
