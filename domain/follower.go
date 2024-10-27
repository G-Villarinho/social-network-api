package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Follower struct {
	ID         uuid.UUID `gorm:"column:id;type:char(36);primaryKey"`
	UserID     uuid.UUID `gorm:"column:userId;type:char(36);not null"`
	FollowerID uuid.UUID `gorm:"column:FollowerID;type:char(36);not null"`
	User       User      `gorm:"foreignKey:UserID"`
	Follower   User      `gorm:"foreignKey:FollowerID"`
	CreatedAt  time.Time `gorm:"column:createdAt;not null"`
	UpdatedAt  time.Time `gorm:"column:updatedAt;default:null"`
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
