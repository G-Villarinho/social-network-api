package domain

//go:generate mockery --name=PostHandler --output=../mocks --outpkg=mocks
//go:generate mockery --name=PostService --output=../mocks --outpkg=mocks
//go:generate mockery --name=PostRepository --output=../mocks --outpkg=mocks

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
	ErrPostNotFound        = errors.New("post not found")
	ErrPostNotBelongToUser = errors.New("post not belong to user")
	ErrPostAlreadyLiked    = errors.New("post already liked")
	ErrPostNotLiked        = errors.New("post not liked")
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
	LikesByUser    bool      `json:"likesByUser"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}

type PostHandler interface {
	CreatePost(ctx echo.Context) error
	GetPosts(ctx echo.Context) error
	GetPostById(ctx echo.Context) error
	UpdatePost(ctx echo.Context) error
	DeletePost(ctx echo.Context) error
	GetByUserID(ctx echo.Context) error
	LikePost(ctx echo.Context) error
	UnlikePost(ctx echo.Context) error
}

type PostService interface {
	CreatePost(ctx context.Context, payload PostPayload) error
	GetPosts(ctx context.Context, page, limit int) (*Pagination[*PostResponse], error)
	GetPostById(ctx context.Context, ID uuid.UUID) (*PostResponse, error)
	UpdatePost(ctx context.Context, ID uuid.UUID, payload PostUpdatePayload) error
	DeletePost(ctx context.Context, ID uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*PostResponse, error)
	LikePost(ctx context.Context, ID uuid.UUID) error
	UnlikePost(ctx context.Context, ID uuid.UUID) error
	ProcessLikePost(ctx context.Context, payload LikePayload) error
	ProcessUnlikePost(ctx context.Context, payload LikePayload) error
}

type PostRepository interface {
	CreatePost(ctx context.Context, post Post) error
	GetPaginatedPosts(ctx context.Context, userID uuid.UUID, page int, limit int) (*Pagination[*Post], error)
	GetPostById(ctx context.Context, ID uuid.UUID, preload bool) (*Post, error)
	UpdatePost(ctx context.Context, ID uuid.UUID, post Post) error
	DeletePost(ctx context.Context, ID uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Post, error)
	LikePost(ctx context.Context, like Like) error
	UnlikePost(ctx context.Context, ID uuid.UUID, userID uuid.UUID) error
	HasUserLikedPost(ctx context.Context, ID uuid.UUID, userID uuid.UUID) (bool, error)
	GetLikedPostIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error)
	GetLikesByPostIDs(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) ([]uuid.UUID, error)
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

	if p.Title == "" && p.Content == "" {
		return ValidationErrors{
			"General": "Title or Content is required",
		}
	}

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

func (pr *PostResponse) SetLikesByUser(likesByUser bool) {
	pr.LikesByUser = likesByUser
}

func (p *Post) Update(payload PostUpdatePayload) {
	if payload.Title != "" {
		p.Title = payload.Title
	}

	if payload.Content != "" {
		p.Content = payload.Content
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
