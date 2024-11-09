// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/G-Villarinho/social-network/domain"
	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// CheckUsername provides a mock function with given fields: ctx, payload
func (_m *UserService) CheckUsername(ctx context.Context, payload domain.CheckUsernamePayload) (*domain.UsernameSuggestionResponse, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for CheckUsername")
	}

	var r0 *domain.UsernameSuggestionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.CheckUsernamePayload) (*domain.UsernameSuggestionResponse, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.CheckUsernamePayload) *domain.UsernameSuggestionResponse); ok {
		r0 = rf(ctx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.UsernameSuggestionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.CheckUsernamePayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, payload
func (_m *UserService) CreateUser(ctx context.Context, payload domain.UserPayload) (string, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserPayload) (string, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserPayload) string); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.UserPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: ctx
func (_m *UserService) DeleteUser(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUser provides a mock function with given fields: ctx
func (_m *UserService) GetUser(ctx context.Context) (*domain.UserResponse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 *domain.UserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*domain.UserResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *domain.UserResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.UserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignIn provides a mock function with given fields: ctx, payload
func (_m *UserService) SignIn(ctx context.Context, payload domain.SignInPayload) (string, error) {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for SignIn")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.SignInPayload) (string, error)); ok {
		return rf(ctx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.SignInPayload) string); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.SignInPayload) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignOut provides a mock function with given fields: ctx
func (_m *UserService) SignOut(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for SignOut")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUser provides a mock function with given fields: ctx, payload
func (_m *UserService) UpdateUser(ctx context.Context, payload domain.UserUpdatePayload) error {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserUpdatePayload) error); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}