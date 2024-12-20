// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/G-Villarinho/social-network/domain"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// SessionRepository is an autogenerated mock type for the SessionRepository type
type SessionRepository struct {
	mock.Mock
}

// CreateSession provides a mock function with given fields: ctx, session
func (_m *SessionRepository) CreateSession(ctx context.Context, session domain.Session) error {
	ret := _m.Called(ctx, session)

	if len(ret) == 0 {
		panic("no return value specified for CreateSession")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Session) error); ok {
		r0 = rf(ctx, session)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSession provides a mock function with given fields: ctx, userId
func (_m *SessionRepository) DeleteSession(ctx context.Context, userId uuid.UUID) error {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSession")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetSessionByUserID provides a mock function with given fields: ctx, userID
func (_m *SessionRepository) GetSessionByUserID(ctx context.Context, userID uuid.UUID) (*domain.Session, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetSessionByUserID")
	}

	var r0 *domain.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*domain.Session, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *domain.Session); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSessionRepository creates a new instance of SessionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSessionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *SessionRepository {
	mock := &SessionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
