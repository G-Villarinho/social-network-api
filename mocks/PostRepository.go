// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/G-Villarinho/social-network/domain"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// PostRepository is an autogenerated mock type for the PostRepository type
type PostRepository struct {
	mock.Mock
}

// CreatePost provides a mock function with given fields: ctx, post
func (_m *PostRepository) CreatePost(ctx context.Context, post domain.Post) error {
	ret := _m.Called(ctx, post)

	if len(ret) == 0 {
		panic("no return value specified for CreatePost")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Post) error); ok {
		r0 = rf(ctx, post)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePost provides a mock function with given fields: ctx, ID
func (_m *PostRepository) DeletePost(ctx context.Context, ID uuid.UUID) error {
	ret := _m.Called(ctx, ID)

	if len(ret) == 0 {
		panic("no return value specified for DeletePost")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, ID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByUserID provides a mock function with given fields: ctx, userID
func (_m *PostRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Post, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserID")
	}

	var r0 []*domain.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]*domain.Post, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []*domain.Post); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLikedPostIDs provides a mock function with given fields: ctx, userID
func (_m *PostRepository) GetLikedPostIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetLikedPostIDs")
	}

	var r0 map[uuid.UUID]bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (map[uuid.UUID]bool, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) map[uuid.UUID]bool); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[uuid.UUID]bool)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLikesByPostIDs provides a mock function with given fields: ctx, userID, postIDs
func (_m *PostRepository) GetLikesByPostIDs(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) ([]uuid.UUID, error) {
	ret := _m.Called(ctx, userID, postIDs)

	if len(ret) == 0 {
		panic("no return value specified for GetLikesByPostIDs")
	}

	var r0 []uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []uuid.UUID) ([]uuid.UUID, error)); ok {
		return rf(ctx, userID, postIDs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []uuid.UUID) []uuid.UUID); ok {
		r0 = rf(ctx, userID, postIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, []uuid.UUID) error); ok {
		r1 = rf(ctx, userID, postIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPaginatedPosts provides a mock function with given fields: ctx, userID, page, limit
func (_m *PostRepository) GetPaginatedPosts(ctx context.Context, userID uuid.UUID, page int, limit int) (*domain.Pagination[*domain.Post], error) {
	ret := _m.Called(ctx, userID, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetPaginatedPosts")
	}

	var r0 *domain.Pagination[*domain.Post]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, int, int) (*domain.Pagination[*domain.Post], error)); ok {
		return rf(ctx, userID, page, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, int, int) *domain.Pagination[*domain.Post]); ok {
		r0 = rf(ctx, userID, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Pagination[*domain.Post])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, int, int) error); ok {
		r1 = rf(ctx, userID, page, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPostById provides a mock function with given fields: ctx, ID, preload
func (_m *PostRepository) GetPostById(ctx context.Context, ID uuid.UUID, preload bool) (*domain.Post, error) {
	ret := _m.Called(ctx, ID, preload)

	if len(ret) == 0 {
		panic("no return value specified for GetPostById")
	}

	var r0 *domain.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, bool) (*domain.Post, error)); ok {
		return rf(ctx, ID, preload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, bool) *domain.Post); ok {
		r0 = rf(ctx, ID, preload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, bool) error); ok {
		r1 = rf(ctx, ID, preload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HasUserLikedPost provides a mock function with given fields: ctx, ID, userID
func (_m *PostRepository) HasUserLikedPost(ctx context.Context, ID uuid.UUID, userID uuid.UUID) (bool, error) {
	ret := _m.Called(ctx, ID, userID)

	if len(ret) == 0 {
		panic("no return value specified for HasUserLikedPost")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (bool, error)); ok {
		return rf(ctx, ID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) bool); ok {
		r0 = rf(ctx, ID, userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, ID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LikePost provides a mock function with given fields: ctx, like
func (_m *PostRepository) LikePost(ctx context.Context, like domain.Like) error {
	ret := _m.Called(ctx, like)

	if len(ret) == 0 {
		panic("no return value specified for LikePost")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Like) error); ok {
		r0 = rf(ctx, like)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnlikePost provides a mock function with given fields: ctx, ID, userID
func (_m *PostRepository) UnlikePost(ctx context.Context, ID uuid.UUID, userID uuid.UUID) error {
	ret := _m.Called(ctx, ID, userID)

	if len(ret) == 0 {
		panic("no return value specified for UnlikePost")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, ID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePost provides a mock function with given fields: ctx, ID, post
func (_m *PostRepository) UpdatePost(ctx context.Context, ID uuid.UUID, post domain.Post) error {
	ret := _m.Called(ctx, ID, post)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePost")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, domain.Post) error); ok {
		r0 = rf(ctx, ID, post)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPostRepository creates a new instance of PostRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPostRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *PostRepository {
	mock := &PostRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
