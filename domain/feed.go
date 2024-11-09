package domain

import "context"

type FeedService interface {
	GenerateFeed(ctx context.Context, page, limit int) (*Pagination[*PostResponse], error)
}
