package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/G-Villarinho/social-network/secure"
	"github.com/G-Villarinho/social-network/utils"
	jsoniter "github.com/json-iterator/go"
)

type userService struct {
	di                *internal.Di
	userRepository    domain.UserRepository
	queueService      domain.QueueService
	clientInfoService domain.ClientInfoService
	sessionService    domain.SessionService
	contextService    domain.ContextService
}

func NewUserService(di *internal.Di) (domain.UserService, error) {
	userRepository, err := internal.Invoke[domain.UserRepository](di)
	if err != nil {
		return nil, err
	}

	sessionService, err := internal.Invoke[domain.SessionService](di)
	if err != nil {
		return nil, err
	}

	contextService, err := internal.Invoke[domain.ContextService](di)
	if err != nil {
		return nil, err
	}

	queueService, err := internal.Invoke[domain.QueueService](di)
	if err != nil {
		return nil, err
	}

	clientInfoService, err := internal.Invoke[domain.ClientInfoService](di)
	if err != nil {
		return nil, err
	}

	return &userService{
		di:                di,
		userRepository:    userRepository,
		sessionService:    sessionService,
		contextService:    contextService,
		queueService:      queueService,
		clientInfoService: clientInfoService,
	}, nil
}

func (u *userService) CreateUser(ctx context.Context, payload domain.UserPayload) (string, error) {
	user, err := u.userRepository.GetUserByUsernameOrEmail(ctx, payload.Username, payload.Email)
	if err != nil {
		return "", fmt.Errorf("error to get user by email: %w", err)
	}

	if user != nil {
		if user.Email == payload.Email {
			return "", domain.ErrEmailAlreadyRegister
		}

		if user.Username == payload.Username {
			return "", domain.ErrUsernameAlreadyExists
		}
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
	user, err := u.userRepository.GetUserByEmailOrUsername(ctx, payload.EmailOrUsername)
	if err != nil {
		return "", fmt.Errorf("error to get user by email or username: %w", err)
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

	go func() {
		clientInfo, err := u.clientInfoService.GetClientInfo(ctx)
		if err != nil {
			slog.Error("get client info", slog.String("error", err.Error()))
			return
		}

		message, err := jsoniter.Marshal(getEmailNotificationTask(user, *clientInfo))
		if err != nil {
			slog.Error("marshal email task event", slog.String("error", err.Error()))
			return
		}

		if err := u.queueService.Publish(domain.QueueSendEmail, message); err != nil {
			slog.Error("publish email event", slog.String("error", err.Error()))
			return
		}
	}()

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

	if payload.Username != "" {
		username, err := u.userRepository.GetUserByUsername(ctx, payload.Username)
		if err != nil {
			return fmt.Errorf("error to get user by username: %w", err)
		}

		if username != nil {
			return domain.ErrUsernameAlreadyExists
		}
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

func (u *userService) DeleteUser(ctx context.Context) error {
	userID := u.contextService.GetUserID(ctx)

	if err := u.userRepository.DeleteUser(ctx, userID); err != nil {
		return err
	}

	if err := u.sessionService.DeleteSession(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (u *userService) CheckUsername(ctx context.Context, payload domain.CheckUsernamePayload) (*domain.UsernameSuggestionResponse, error) {
	exits, err := u.userRepository.CheckUsername(ctx, payload.Username)
	if err != nil {
		return nil, fmt.Errorf("check username: %w", err)
	}

	if exits {
		return &domain.UsernameSuggestionResponse{
			Suggestions: utils.GenerateSuggestions(payload.Username, 3),
		}, domain.ErrUsernameAlreadyExists
	}

	return nil, nil
}

func getEmailNotificationTask(user *domain.User, clientInfo domain.ClientInfoResponse) domain.EmailPayloadTask {
	return domain.EmailPayloadTask{
		Template: domain.SignInNotification,
		Subject:  "New Sign-In Detected",
		Recipient: domain.Recipient{
			Name:  fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			Email: user.Email,
		},
		Params: map[string]string{
			"name":      fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			"device":    clientInfo.Device,
			"location":  clientInfo.Location,
			"date_time": clientInfo.LoginTime,
		},
	}
}
