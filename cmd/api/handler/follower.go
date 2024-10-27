package handler

import (
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type followerHandler struct {
	di              *pkg.Di
	followerService domain.FollowerService
}

func NewFollowerHandler(di *pkg.Di) (domain.FollowerHandler, error) {
	followerService, err := pkg.Invoke[domain.FollowerService](di)
	if err != nil {
		return nil, err
	}

	return &followerHandler{
		di:              di,
		followerService: followerService,
	}, nil
}

func (f *followerHandler) FollowUser(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "FollowUser"),
	)

	followerID, err := uuid.Parse(ctx.Param("followerId"))
	if err != nil {
		log.Warn("Error to parse UUID", slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid UUID", "The ID provided is not a valid UUID.")
	}

	if err := f.followerService.FollowUser(ctx.Request().Context(), followerID); err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrFollowerNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The follower does not exist.")
		}

		if err == domain.ErrUserCannotFollowItself {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, nil, "Conflict", "The user cannot follow itself.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusNoContent)
}
