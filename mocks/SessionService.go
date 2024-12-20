// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/G-Villarinho/social-network/domain"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// SessionService is an autogenerated mock type for the SessionService type
type SessionService struct {
	mock.Mock
}

// CreateSession provides a mock function with given fields: ctx, user
func (_m *SessionService) CreateSession(ctx context.Context, user domain.User) (string, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateSession")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.User) (string, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.User) string); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSession provides a mock function with given fields: ctx, userID
func (_m *SessionService) DeleteSession(ctx context.Context, userID uuid.UUID) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSession")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetSessionByToken provides a mock function with given fields: ctx, token
func (_m *SessionService) GetSessionByToken(ctx context.Context, token string) (*domain.Session, error) {
	ret := _m.Called(ctx, token)

	if len(ret) == 0 {
		panic("no return value specified for GetSessionByToken")
	}

	var r0 *domain.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.Session, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Session); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSessionService creates a new instance of SessionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSessionService(t interface {
	mock.TestingT
	Cleanup(func())
}) *SessionService {
	mock := &SessionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
