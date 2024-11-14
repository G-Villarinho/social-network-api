package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/labstack/echo/v4"
)

type feedHandler struct {
	di          *internal.Di
	feedService domain.FeedService
}

func NewFeedHandler(di *internal.Di) (domain.FeedHandler, error) {
	feedService, err := internal.Invoke[domain.FeedService](di)
	if err != nil {
		return nil, err
	}

	return &feedHandler{
		di:          di,
		feedService: feedService,
	}, nil
}

func (f *feedHandler) GetFeed(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "feed"),
		slog.String("func", "GetFeed"),
	)

	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	response, err := f.feedService.GetFeed(ctx.Request().Context(), page, limit)
	if err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrPostNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "No posts found.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}
