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

type statusType string

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrEmailAlreadyRegister     = errors.New("email already exists")
	ErrInvalidPassword          = errors.New("invalid password")
	ErrUserNotFoundInContext    = errors.New("user not found in context")
	ErrHashExpired              = errors.New("the 2FA hash has expired")
	ErrCodeOTPExpired           = errors.New("the code OTP has expired")
	ErrCodeOTPWrong             = errors.New("the code OTP is wrong")
	ErrEmailConfirmationPending = errors.New("email confirmation is pending")
)

const (
	Active   statusType = "active"
	Inactive statusType = "inactive"
	Block    statusType = "block"
)

type User struct {
	ID        uuid.UUID  `gorm:"column:id;type:char(36);primaryKey"`
	FirstName string     `gorm:"column:firstName;type:varchar(255);not null"`
	LastName  string     `gorm:"column:lastName;type:varchar(255);not null"`
	Email     string     `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	Password  string     `gorm:"column:password;type:varchar(255);not null"`
	Avatar    string     `gorm:"column:avatar;type:varchar(255);default:null"`
	Status    statusType `gorm:"type:enum('active','inactive','block');default:'active';index"`
	CreatedAt time.Time  `gorm:"column:createdAt;not null"`
	UpdatedAt time.Time  `gorm:"column:updatedAt;default:null"`
}

type UserPayload struct {
	FirstName string `json:"firstName" validate:"required,min=1,max=255"`
	LastName  string `json:"lastName" validate:"required,min=1,max=255"`
	Password  string `json:"password" validate:"required,strongpassword"`
	Email     string `json:"email" validate:"required,email,max=255"`
}

type UserUpdatePayload struct {
	FirstName string `json:"firstName" validate:"omitempty,min=1,max=255"`
	LastName  string `json:"lastName" validate:"omitempty,min=1,max=255"`
}

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
}

type SignInPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserHandler interface {
	CreateUser(ctx echo.Context) error
	SignIn(ctx echo.Context) error
	SignOut(ctx echo.Context) error
	GetUser(ctx echo.Context) error
	UpdateUser(ctx echo.Context) error
}

type UserService interface {
	CreateUser(ctx context.Context, payload UserPayload) (string, error)
	SignIn(ctx context.Context, payload SignInPayload) (string, error)
	SignOut(ctx context.Context) error
	GetUser(ctx context.Context) (*UserResponse, error)
	UpdateUser(ctx context.Context, payload UserUpdatePayload) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, user User) error
}

func (u *UserPayload) trim() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
}

func (s *SignInPayload) trim() {
	s.Email = strings.TrimSpace(strings.ToLower(s.Email))
}

func (uup *UserUpdatePayload) trim() {
	uup.FirstName = strings.TrimSpace(uup.FirstName)
	uup.LastName = strings.TrimSpace(uup.LastName)
}

func (u *UserPayload) Validate() ValidationErrors {
	u.trim()
	return ValidateStruct(u)
}

func (s *SignInPayload) Validate() ValidationErrors {
	s.trim()
	return ValidateStruct(s)
}

func (uup *UserUpdatePayload) Validate() ValidationErrors {
	uup.trim()

	if uup.FirstName == "" && uup.LastName == "" {
		return ValidationErrors{"General": "firstName or lastName is required"}
	}

	return ValidateStruct(uup)
}

func (User) TableName() string {
	return "User"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	u.CreatedAt = time.Now().UTC()
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now().UTC()
	return
}

func (up *UserPayload) ToUser(passwordHash string) *User {
	return &User{
		FirstName: up.FirstName,
		LastName:  up.LastName,
		Email:     up.Email,
		Password:  passwordHash,
	}
}

func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Avatar:    u.Avatar,
	}
}

func (u *User) Update(payload UserUpdatePayload) {
	if payload.FirstName != "" {
		u.FirstName = payload.FirstName
	}

	if payload.LastName != "" {
		u.LastName = payload.LastName
	}
}
