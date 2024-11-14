package domain

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
)

var (
	ErrFeedNotFound = errors.New("feed not found")
)

type FeedHandler interface {
	GetFeed(ctx echo.Context) error
}

type FeedService interface {
	GetFeed(ctx context.Context, page, limit int) (*Pagination[*PostResponse], error)
}
