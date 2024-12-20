// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/G-Villarinho/social-network/domain"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MemoryCacheRepository is an autogenerated mock type for the MemoryCacheRepository type
type MemoryCacheRepository struct {
	mock.Mock
}

// GetCachedLikes provides a mock function with given fields: ctx, userID, postIDs
func (_m *MemoryCacheRepository) GetCachedLikes(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) (*domain.LikeCache, error) {
	ret := _m.Called(ctx, userID, postIDs)

	if len(ret) == 0 {
		panic("no return value specified for GetCachedLikes")
	}

	var r0 *domain.LikeCache
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []uuid.UUID) (*domain.LikeCache, error)); ok {
		return rf(ctx, userID, postIDs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []uuid.UUID) *domain.LikeCache); ok {
		r0 = rf(ctx, userID, postIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.LikeCache)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, []uuid.UUID) error); ok {
		r1 = rf(ctx, userID, postIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPosts provides a mock function with given fields: ctx, userID, page, limit
func (_m *MemoryCacheRepository) GetPosts(ctx context.Context, userID uuid.UUID, page int, limit int) (*domain.Pagination[*domain.PostResponse], error) {
	ret := _m.Called(ctx, userID, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetPosts")
	}

	var r0 *domain.Pagination[*domain.PostResponse]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, int, int) (*domain.Pagination[*domain.PostResponse], error)); ok {
		return rf(ctx, userID, page, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, int, int) *domain.Pagination[*domain.PostResponse]); ok {
		r0 = rf(ctx, userID, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Pagination[*domain.PostResponse])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, int, int) error); ok {
		r1 = rf(ctx, userID, page, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemovePostLike provides a mock function with given fields: ctx, postID, userID
func (_m *MemoryCacheRepository) RemovePostLike(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	ret := _m.Called(ctx, postID, userID)

	if len(ret) == 0 {
		panic("no return value specified for RemovePostLike")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, postID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetLikesByPostIDs provides a mock function with given fields: ctx, userID, postIDs
func (_m *MemoryCacheRepository) SetLikesByPostIDs(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) error {
	ret := _m.Called(ctx, userID, postIDs)

	if len(ret) == 0 {
		panic("no return value specified for SetLikesByPostIDs")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []uuid.UUID) error); ok {
		r0 = rf(ctx, userID, postIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetPost provides a mock function with given fields: ctx, userID, posts, page, limit
func (_m *MemoryCacheRepository) SetPost(ctx context.Context, userID uuid.UUID, posts *domain.Pagination[*domain.PostResponse], page int, limit int) error {
	ret := _m.Called(ctx, userID, posts, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for SetPost")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, *domain.Pagination[*domain.PostResponse], int, int) error); ok {
		r0 = rf(ctx, userID, posts, page, limit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetPostLike provides a mock function with given fields: ctx, postID, userID
func (_m *MemoryCacheRepository) SetPostLike(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	ret := _m.Called(ctx, postID, userID)

	if len(ret) == 0 {
		panic("no return value specified for SetPostLike")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, postID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMemoryCacheRepository creates a new instance of MemoryCacheRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMemoryCacheRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MemoryCacheRepository {
	mock := &MemoryCacheRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
