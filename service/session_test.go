package service

import (
	"context"
	"errors"
	"testing"

	"github.com/G-Villarinho/social-network/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSession_WhenSessionDeletedSuccessfully_ShouldNotReturnError(t *testing.T) {
	ctx := context.Background()
	sessionRepoMock := new(mocks.SessionRepository)

	userID := uuid.New()
	sessionService := &sessionService{
		sessionRepository: sessionRepoMock,
	}

	sessionRepoMock.On("DeleteSession", ctx, userID).Return(nil)

	err := sessionService.DeleteSession(ctx, userID)

	assert.NoError(t, err)
	sessionRepoMock.AssertExpectations(t)
}

func TestDeleteSession_WhenRepositoryFails_ShouldReturnError(t *testing.T) {
	ctx := context.Background()
	sessionRepoMock := new(mocks.SessionRepository)

	userID := uuid.New()
	sessionService := &sessionService{
		sessionRepository: sessionRepoMock,
	}

	sessionRepoMock.On("DeleteSession", ctx, userID).Return(errors.New("repository error"))

	err := sessionService.DeleteSession(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository error")
	sessionRepoMock.AssertExpectations(t)
}

func TestExtractSessionFromToken_WhenTokenIsInvalid_ShouldReturnError(t *testing.T) {
	invalidToken := "invalid-token"

	sessionService := &sessionService{}

	session, err := sessionService.extractSessionFromToken(invalidToken)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error to parse token")
	assert.Nil(t, session)
}
