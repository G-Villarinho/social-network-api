// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/G-Villarinho/social-network/domain"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// FollowerRepository is an autogenerated mock type for the FollowerRepository type
type FollowerRepository struct {
	mock.Mock
}

// CreateFollower provides a mock function with given fields: ctx, follower
func (_m *FollowerRepository) CreateFollower(ctx context.Context, follower domain.Follower) error {
	ret := _m.Called(ctx, follower)

	if len(ret) == 0 {
		panic("no return value specified for CreateFollower")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Follower) error); ok {
		r0 = rf(ctx, follower)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteFollower provides a mock function with given fields: ctx, followerId
func (_m *FollowerRepository) DeleteFollower(ctx context.Context, followerId uuid.UUID) error {
	ret := _m.Called(ctx, followerId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteFollower")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, followerId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetFollower provides a mock function with given fields: ctx, userID, followerId
func (_m *FollowerRepository) GetFollower(ctx context.Context, userID uuid.UUID, followerId uuid.UUID) (*domain.Follower, error) {
	ret := _m.Called(ctx, userID, followerId)

	if len(ret) == 0 {
		panic("no return value specified for GetFollower")
	}

	var r0 *domain.Follower
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (*domain.Follower, error)); ok {
		return rf(ctx, userID, followerId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) *domain.Follower); ok {
		r0 = rf(ctx, userID, followerId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Follower)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, userID, followerId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFollowers provides a mock function with given fields: ctx, userID
func (_m *FollowerRepository) GetFollowers(ctx context.Context, userID uuid.UUID) ([]*domain.Follower, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetFollowers")
	}

	var r0 []*domain.Follower
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]*domain.Follower, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []*domain.Follower); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Follower)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFollowings provides a mock function with given fields: ctx, userID
func (_m *FollowerRepository) GetFollowings(ctx context.Context, userID uuid.UUID) ([]*domain.Follower, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetFollowings")
	}

	var r0 []*domain.Follower
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]*domain.Follower, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []*domain.Follower); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Follower)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFollowerRepository creates a new instance of FollowerRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFollowerRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *FollowerRepository {
	mock := &FollowerRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}