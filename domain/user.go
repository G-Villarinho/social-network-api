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

//go:generate mockery --name=UserHandler --output=../mocks --outpkg=mocks
//go:generate mockery --name=UserService --output=../mocks --outpkg=mocks
//go:generate mockery --name=UserRepository --output=../mocks --outpkg=mocks
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
	ErrUsernameAlreadyExists    = errors.New("username already exists")
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
	Username  string     `gorm:"column:username;type:varchar(20);uniqueIndex;not null"`
	Email     string     `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	Password  string     `gorm:"column:password;type:varchar(255);not null"`
	Avatar    string     `gorm:"column:avatar;type:varchar(255);default:null"`
	Status    statusType `gorm:"type:enum('active','inactive','block');default:'active';index"`
	CreatedAt time.Time  `gorm:"column:createdAt;not null"`
	UpdatedAt time.Time  `gorm:"column:updatedAt;default:null"`
}

type UserPayload struct {
	FirstName       string `json:"firstName" validate:"required,max=255"`
	LastName        string `json:"lastName" validate:"required,max=255"`
	Email           string `json:"email" validate:"required,email,max=255"`
	Username        string `json:"username" validate:"required,username,min=3,max=20"`
	Password        string `json:"password" validate:"required,strongpassword"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type UserUpdatePayload struct {
	FirstName string `json:"firstName" validate:"omitempty,min=1,max=255"`
	LastName  string `json:"lastName" validate:"omitempty,min=1,max=255"`
	Username  string `json:"username" validate:"omitempty,username,min=3,max=20"`
}

type CheckUsernamePayload struct {
	Username string `json:"username" validate:"required,username"`
}

type CheckPasswordStrongPayload struct {
	Password string `json:"password" validate:"required,strongpassword"`
}

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
}

type UserFollowerResponse struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
}

type UsernameSuggestionResponse struct {
	Suggestions []string `json:"suggestions"`
}

type SignInPayload struct {
	EmailOrUsername string `json:"emailOrUsername" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
}

type UserHandler interface {
	CreateUser(ctx echo.Context) error
	SignIn(ctx echo.Context) error
	SignOut(ctx echo.Context) error
	GetUser(ctx echo.Context) error
	UpdateUser(ctx echo.Context) error
	DeleteUser(ctx echo.Context) error
	CheckUsername(ctx echo.Context) error
	CheckPasswordStrong(ctx echo.Context) error
}

type UserService interface {
	CreateUser(ctx context.Context, payload UserPayload) (string, error)
	SignIn(ctx context.Context, payload SignInPayload) (string, error)
	SignOut(ctx context.Context) error
	GetUser(ctx context.Context) (*UserResponse, error)
	UpdateUser(ctx context.Context, payload UserUpdatePayload) error
	DeleteUser(ctx context.Context) error
	CheckUsername(ctx context.Context, payload CheckUsernamePayload) (*UsernameSuggestionResponse, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, ID uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, user User) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*User, error)
	DeleteUser(ctx context.Context, ID uuid.UUID) error
	GetUserByUsernameOrEmail(ctx context.Context, username, email string) (*User, error)
	CheckUsername(ctx context.Context, username string) (bool, error)
}

func (u *UserPayload) trim() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
}

func (s *SignInPayload) trim() {
	s.EmailOrUsername = strings.TrimSpace(strings.ToLower(s.EmailOrUsername))
}

func (uup *UserUpdatePayload) trim() {
	uup.FirstName = strings.TrimSpace(uup.FirstName)
	uup.LastName = strings.TrimSpace(uup.LastName)
}

func (c *CheckUsernamePayload) trim() {
	c.Username = strings.TrimSpace(c.Username)
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

	if uup.FirstName == "" && uup.LastName == "" && uup.Username == "" {
		return ValidationErrors{"General": "firstName or lastName or username is required"}
	}

	return ValidateStruct(uup)
}

func (c *CheckUsernamePayload) Validate() ValidationErrors {
	c.trim()
	return ValidateStruct(c)
}

func (c *CheckPasswordStrongPayload) Validate() ValidationErrors {
	return ValidateStruct(c)
}

func (u *User) ToUserFollowerResponse() *UserFollowerResponse {
	return &UserFollowerResponse{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
	}
}

func (up *UserPayload) ToUser(passwordHash string) *User {
	return &User{
		ID:        uuid.New(),
		FirstName: up.FirstName,
		LastName:  up.LastName,
		Username:  up.Username,
		Email:     up.Email,
		Password:  passwordHash,
	}
}

func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
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

	if payload.Username != "" {
		u.Username = payload.Username
	}
}

func (User) TableName() string {
	return "User"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now().UTC()
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now().UTC()
	return
}
