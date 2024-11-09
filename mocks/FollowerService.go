// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/G-Villarinho/social-network/domain"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// FollowerService is an autogenerated mock type for the FollowerService type
type FollowerService struct {
	mock.Mock
}

// FollowUser provides a mock function with given fields: ctx, userID
func (_m *FollowerService) FollowUser(ctx context.Context, userID uuid.UUID) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for FollowUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetFollowers provides a mock function with given fields: ctx
func (_m *FollowerService) GetFollowers(ctx context.Context) ([]*domain.FollowerResponse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetFollowers")
	}

	var r0 []*domain.FollowerResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*domain.FollowerResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*domain.FollowerResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.FollowerResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFollowings provides a mock function with given fields: ctx
func (_m *FollowerService) GetFollowings(ctx context.Context) ([]*domain.FollowerResponse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetFollowings")
	}

	var r0 []*domain.FollowerResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*domain.FollowerResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*domain.FollowerResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.FollowerResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnfollowUser provides a mock function with given fields: ctx, userID
func (_m *FollowerService) UnfollowUser(ctx context.Context, userID uuid.UUID) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for UnfollowUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFollowerService creates a new instance of FollowerService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFollowerService(t interface {
	mock.TestingT
	Cleanup(func())
}) *FollowerService {
	mock := &FollowerService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}