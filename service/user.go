package service

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/G-Villarinho/social-network/secure"
)

type userService struct {
	di             *pkg.Di
	userRepository domain.UserRepository
	sessionService domain.SessionService
}

func NewUserService(di *pkg.Di) (domain.UserService, error) {
	userRepository, err := pkg.Invoke[domain.UserRepository](di)
	if err != nil {
		return nil, err
	}

	sessionService, err := pkg.Invoke[domain.SessionService](di)
	if err != nil {
		return nil, err
	}

	return &userService{
		di:             di,
		userRepository: userRepository,
		sessionService: sessionService,
	}, nil
}

func (u *userService) CreateUser(ctx context.Context, payload domain.UserPayload) (string, error) {
	user, err := u.userRepository.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return "", err
	}

	if user != nil {
		return "", domain.ErrEmailAlreadyRegister
	}

	passwordHash, err := secure.HashPassword(payload.Password)
	if err != nil {
		return "", fmt.Errorf("error to hash password: %w", err)
	}

	user = payload.ToUser(string(passwordHash))
	if err := u.userRepository.CreateUser(ctx, *user); err != nil {
		return "", err
	}

	token, err := u.sessionService.CreateSession(ctx, *user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *userService) SignIn(ctx context.Context, payload domain.SignInPayload) (string, error) {
	user, err := u.userRepository.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return "", fmt.Errorf("error to get user by email: %w", err)
	}

	if user == nil {
		return "", domain.ErrUserNotFound
	}

	if err := secure.CheckPassword(user.Password, payload.Password); err != nil {
		return "", domain.ErrInvalidPassword
	}

	token, err := u.sessionService.CreateSession(ctx, *user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *userService) SignOut(ctx context.Context) error {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return domain.ErrSessionNotFound
	}

	if err := u.sessionService.DeleteSession(ctx, session.UserID); err != nil {
		return err
	}

	return nil
}

func (u *userService) GetUser(ctx context.Context) (*domain.UserResponse, error) {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return nil, domain.ErrSessionNotFound
	}

	user, err := u.userRepository.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return user.ToUserResponse(), nil
}

func (u *userService) UpdateUser(ctx context.Context, payload domain.UserUpdatePayload) error {
	session, ok := ctx.Value(domain.SessionKey).(*domain.Session)
	if !ok {
		return domain.ErrSessionNotFound
	}

	user, err := u.userRepository.GetUserByID(ctx, session.UserID)
	if err != nil {
		return fmt.Errorf("error to get user by ID: %w", err)
	}

	if user == nil {
		return domain.ErrUserNotFound
	}

	user.Update(payload)

	if err := u.userRepository.UpdateUser(ctx, *user); err != nil {
		return err
	}

	return nil
}
