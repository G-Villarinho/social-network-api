package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_WhenUserAlreadyExists_ShouldReturnErrorAlreadyRegister(t *testing.T) {
	ctx := context.Background()

	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	existingUser := &domain.User{Email: "gabriel@test.com", Username: "gabriel"}

	payload := domain.UserPayload{
		FirstName: "gabriel",
		LastName:  "Soares",
		Email:     "gabriel@test.com",
		Username:  "gabriel",
		Password:  "password123",
	}

	userRepoMock.On("GetUserByUsernameOrEmail", ctx, payload.Username, payload.Email).Return(existingUser, nil)

	token, err := userService.CreateUser(ctx, payload)
	assert.Equal(t, domain.ErrEmailAlreadyRegister, err)
	assert.Empty(t, token)
	userRepoMock.AssertExpectations(t)
}

func TestCreateUser_WhenGetUserByUsernameOrEmailFails_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	payload := domain.UserPayload{
		FirstName: "Gabriel",
		LastName:  "Soares",
		Email:     "gabriel@test.com",
		Username:  "gabriel",
		Password:  "password123",
	}

	userRepoMock.On("GetUserByUsernameOrEmail", ctx, payload.Username, payload.Email).Return(nil, errors.New("repository error"))

	token, err := userService.CreateUser(ctx, payload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository error")
	assert.Empty(t, token)
	userRepoMock.AssertExpectations(t)
}

func TestCreateUser_WhenCreateUserRepositoryFails_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	payload := domain.UserPayload{
		FirstName: "Gabriel",
		LastName:  "Soares",
		Email:     "gabriel@test.com",
		Username:  "gabriel",
		Password:  "password123",
	}

	userRepoMock.On("GetUserByUsernameOrEmail", ctx, payload.Username, payload.Email).Return(nil, nil)
	userRepoMock.On("CreateUser", ctx, mock.Anything).Return(errors.New("repository error"))

	token, err := userService.CreateUser(ctx, payload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository error")
	assert.Empty(t, token)
	userRepoMock.AssertExpectations(t)
}

func TestCreateUser_WhenUsernameAlreadyExists_ShouldReturnErrorUsernameAlreadyExists(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	existingUser := &domain.User{
		Email:    "gabriel@teste.com",
		Username: "gabriel",
	}

	payload := domain.UserPayload{
		FirstName: "Gabriel",
		LastName:  "Soares",
		Email:     "gabriel01@test.com",
		Username:  "gabriel",
		Password:  "password123",
	}

	userRepoMock.On("GetUserByUsernameOrEmail", ctx, payload.Username, payload.Email).Return(existingUser, nil)

	token, err := userService.CreateUser(ctx, payload)

	assert.Equal(t, domain.ErrUsernameAlreadyExists, err)
	assert.Empty(t, token)
	userRepoMock.AssertExpectations(t)
}

func TestCreateUser_WhenHashPasswordFails_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	payload := domain.UserPayload{
		FirstName: "Gabriel",
		LastName:  "Villarinho",
		Email:     "gabriel@test.com",
		Username:  "gabriel",
		Password:  strings.Repeat("a", 1000),
	}

	userRepoMock.On("GetUserByUsernameOrEmail", ctx, payload.Username, payload.Email).Return(nil, nil)
	userRepoMock.On("CreateUser", ctx, mock.Anything).Return(nil)

	token, err := userService.CreateUser(ctx, payload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hash")
	assert.Empty(t, token)
	userRepoMock.AssertNotCalled(t, "CreateUser", ctx, mock.Anything)
}

func TestSignIn_WhenUserNotFound_ShouldReturnErrorUserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	payload := domain.SignInPayload{
		EmailOrUsername: "gabriel",
		Password:        "password123",
	}

	userRepoMock.On("GetUserByEmailOrUsername", ctx, payload.EmailOrUsername).Return(nil, nil)

	token, err := userService.SignIn(ctx, payload)

	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Empty(t, token)
	userRepoMock.AssertExpectations(t)
}

func TestSignIn_WhenInvalidPassword_ShouldReturnErrorInvalidPassword(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	user := &domain.User{
		ID:       uuid.New(),
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "$2a$10$rjs5yVRXcCvjCdF1zRyHTu3wtsRXlVjP/YXJ0BzqCzYrMM2w7UjJG",
	}
	payload := domain.SignInPayload{
		EmailOrUsername: "gabriel",
		Password:        "wrong_password",
	}

	userRepoMock.On("GetUserByEmailOrUsername", ctx, payload.EmailOrUsername).Return(user, nil)
	sessionServiceMock.On("CreateSession", ctx, *user).Return("", nil)

	token, err := userService.SignIn(ctx, payload)

	assert.Equal(t, domain.ErrInvalidPassword, err)
	assert.Empty(t, token)
	userRepoMock.AssertExpectations(t)
}

func TestSignIn_WhenSuccessful_ShouldReturnToken(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	user := &domain.User{
		ID:       uuid.New(),
		Username: "gabriel",
		Email:    "gabriel@test.com",
		Password: "$2a$10$rjs5yVRXcCvjCdF1zRyHTu3wtsRXlVjP/YXJ0BzqCzYrMM2w7UjJG",
	}
	payload := domain.SignInPayload{
		EmailOrUsername: "gabriel",
		Password:        "Abc@123456",
	}

	userRepoMock.On("GetUserByEmailOrUsername", ctx, payload.EmailOrUsername).Return(user, nil)
	sessionServiceMock.On("CreateSession", ctx, *user).Return("valid-token", nil)

	token, err := userService.SignIn(ctx, payload)

	assert.NoError(t, err)
	assert.Equal(t, "valid-token", token)
	userRepoMock.AssertExpectations(t)
	sessionServiceMock.AssertExpectations(t)
}

func TestSignOut_WhenSessionNotFound_ShouldReturnErrorSessionNotFound(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	err := userService.SignOut(ctx)

	assert.Equal(t, domain.ErrSessionNotFound, err)
}

func TestSignOut_WhenSuccessful_ShouldCompleteWithoutError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)
	sessionServiceMock.On("DeleteSession", ctx, session.UserID).Return(nil)

	err := userService.SignOut(ctx)

	assert.NoError(t, err)
	sessionServiceMock.AssertExpectations(t)
}

func TestGetUser_WhenSessionNotFound_ShouldReturnErrorSessionNotFound(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	user, err := userService.GetUser(ctx)

	assert.Equal(t, domain.ErrSessionNotFound, err)
	assert.Nil(t, user)
}

func TestGetUser_WhenUserNotFound_ShouldReturnErrorUserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)
	userRepoMock.On("GetUserByID", ctx, session.UserID).Return(nil, nil)

	user, err := userService.GetUser(ctx)

	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, user)
	userRepoMock.AssertExpectations(t)
}

func TestGetUser_WhenSuccessful_ShouldReturnUser(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	user := &domain.User{
		ID:        session.UserID,
		FirstName: "gabriel",
		LastName:  "Villarinho",
		Email:     "gabriel@test.com",
	}

	ctx = context.WithValue(ctx, domain.SessionKey, session)
	userRepoMock.On("GetUserByID", ctx, session.UserID).Return(user, nil)

	userResponse, err := userService.GetUser(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, userResponse)
	assert.Equal(t, user.FirstName, userResponse.FirstName)
	assert.Equal(t, user.LastName, userResponse.LastName)
	assert.Equal(t, user.Email, userResponse.Email)
	userRepoMock.AssertExpectations(t)
}

func TestUpdateUser_WhenSessionNotFound_ShouldReturnErrorSessionNotFound(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	payload := domain.UserUpdatePayload{
		FirstName: "Gabriel",
	}

	err := userService.UpdateUser(ctx, payload)

	assert.Equal(t, domain.ErrSessionNotFound, err)
}

func TestUpdateUser_WhenUsernameAlreadyExists_ShouldReturnErrorUsernameAlreadyExists(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	payload := domain.UserUpdatePayload{
		Username: "gabriel",
	}
	existingUser := &domain.User{ID: uuid.New(), Username: "gabriel"}

	userRepoMock.On("GetUserByID", ctx, session.UserID).Return(&domain.User{ID: session.UserID}, nil)

	userRepoMock.On("GetUserByUsername", ctx, payload.Username).Return(existingUser, nil)

	err := userService.UpdateUser(ctx, payload)

	assert.Equal(t, domain.ErrUsernameAlreadyExists, err)
	userRepoMock.AssertExpectations(t)
}

func TestUpdateUser_WhenUserNotFound_ShouldReturnErrorUserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	payload := domain.UserUpdatePayload{
		FirstName: "Gabriel",
	}

	userRepoMock.On("GetUserByID", ctx, session.UserID).Return(nil, nil)

	err := userService.UpdateUser(ctx, payload)

	assert.Equal(t, domain.ErrUserNotFound, err)
	userRepoMock.AssertExpectations(t)
}

func TestUpdateUser_WhenSuccessful_ShouldUpdateUser(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	session := &domain.Session{UserID: uuid.New()}
	ctx = context.WithValue(ctx, domain.SessionKey, session)

	user := &domain.User{ID: session.UserID, FirstName: "Gabiel"}
	payload := domain.UserUpdatePayload{
		FirstName: "Gabriel",
	}

	userRepoMock.On("GetUserByID", ctx, session.UserID).Return(user, nil)

	updatedUser := *user
	updatedUser.FirstName = "Gabriel"
	userRepoMock.On("UpdateUser", ctx, updatedUser).Return(nil)

	err := userService.UpdateUser(ctx, payload)

	assert.NoError(t, err)
	assert.Equal(t, "Gabriel", user.FirstName)
	userRepoMock.AssertExpectations(t)
}

func TestDeleteUser_WhenUserNotFound_ShouldReturnErrorUserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	userID := uuid.New()
	contextServiceMock.On("GetUserID", ctx).Return(userID)
	userRepoMock.On("DeleteUser", ctx, userID).Return(domain.ErrUserNotFound)

	err := userService.DeleteUser(ctx)

	assert.Equal(t, domain.ErrUserNotFound, err)
	userRepoMock.AssertExpectations(t)
}

func TestDeleteUser_WhenSuccessful_ShouldCompleteWithoutError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	userID := uuid.New()
	contextServiceMock.On("GetUserID", ctx).Return(userID)
	userRepoMock.On("DeleteUser", ctx, userID).Return(nil)
	sessionServiceMock.On("DeleteSession", ctx, userID).Return(nil)

	err := userService.DeleteUser(ctx)

	assert.NoError(t, err)
	userRepoMock.AssertExpectations(t)
	sessionServiceMock.AssertExpectations(t)
}

func TestCheckUsername_WhenUsernameExists_ShouldReturnSuggestions(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	payload := domain.CheckUsernamePayload{
		Username: "gabriel",
	}
	userRepoMock.On("CheckUsername", ctx, payload.Username).Return(true, nil)

	response, err := userService.CheckUsername(ctx, payload)

	assert.Equal(t, domain.ErrUsernameAlreadyExists, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Suggestions)
	userRepoMock.AssertExpectations(t)
}

func TestCheckUsername_WhenUsernameAvailable_ShouldCompleteWithoutError(t *testing.T) {
	ctx := context.Background()
	userRepoMock := new(mocks.UserRepository)
	sessionServiceMock := new(mocks.SessionService)
	contextServiceMock := new(mocks.ContextService)

	userService := &userService{
		userRepository: userRepoMock,
		sessionService: sessionServiceMock,
		contextService: contextServiceMock,
	}

	payload := domain.CheckUsernamePayload{
		Username: "gabriel",
	}
	userRepoMock.On("CheckUsername", ctx, payload.Username).Return(false, nil)

	response, err := userService.CheckUsername(ctx, payload)

	assert.NoError(t, err)
	assert.Nil(t, response)
	userRepoMock.AssertExpectations(t)
}
