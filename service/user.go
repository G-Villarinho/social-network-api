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