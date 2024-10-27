package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type Post struct {
	ID        uuid.UUID `gorm:"column:id;type:char(36);primaryKey"`
	AuthorID  uuid.UUID `gorm:"column:authorID;type:char(36);not null"`
	Author    User      `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	Likes     uint64    `gorm:"column:likes;not null;default:0"`
	Title     string    `gorm:"column:title;type:varchar(50);not null"`
	Content   string    `gorm:"column:content;type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"column:createdAt;not null"`
	UpdatedAt time.Time `gorm:"column:updatedAt;default:null"`
}

type PostPayload struct {
	Title   string `json:"title" validate:"required,max=50"`
	Content string `json:"content" validate:"required,max=255"`
}

type PostUpdatePayload struct {
	Title   string `json:"title" validate:"omitempty,max=50"`
	Content string `json:"content" validate:"omitempty,max=255"`
}

type PostResponse struct {
	ID             uuid.UUID `json:"id"`
	AuthorUsername string    `json:"authorUsername"`
	Likes          uint64    `json:"likes"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}

type PostHandler interface {
	CreatePost(ctx echo.Context) error
	GetPosts(ctx echo.Context) error
}

type PostService interface {
	CreatePost(ctx context.Context, payload PostPayload) error
	GetPosts(ctx context.Context) ([]*PostResponse, error)
}

type PostRepository interface {
	CreatePost(ctx context.Context, post Post) error
	GetPosts(ctx context.Context, userID uuid.UUID) ([]*Post, error)
}

func (p *PostPayload) trim() {
	p.Title = strings.TrimSpace(p.Title)
	p.Content = strings.TrimSpace(p.Content)
}

func (p *PostUpdatePayload) trim() {
	p.Title = strings.TrimSpace(p.Title)
	p.Content = strings.TrimSpace(p.Content)
}

func (p *PostPayload) Validate() ValidationErrors {
	p.trim()
	return ValidateStruct(p)
}

func (p *PostUpdatePayload) Validate() ValidationErrors {
	p.trim()
	return ValidateStruct(p)
}

func (p *PostPayload) ToPost(userId uuid.UUID) *Post {
	return &Post{
		AuthorID: userId,
		Title:    p.Title,
		Content:  p.Content,
	}
}

func (p *Post) ToPostResponse() *PostResponse {
	return &PostResponse{
		ID:             p.ID,
		AuthorUsername: p.Author.Username,
		Likes:          p.Likes,
		Title:          p.Title,
		Content:        p.Content,
		CreatedAt:      p.CreatedAt,
	}
}

func (Post) TableName() string {
	return "Post"
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	p.CreatedAt = time.Now().UTC()
	return
}

func (p *Post) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now().UTC()
	return
}
