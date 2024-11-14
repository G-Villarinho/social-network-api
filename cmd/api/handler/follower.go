package handler

import (
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

type followerHandler struct {
	di              *internal.Di
	followerService domain.FollowerService
}

func NewFollowerHandler(di *internal.Di) (domain.FollowerHandler, error) {
	followerService, err := internal.Invoke[domain.FollowerService](di)
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
		slog.String("handler", "follower"),
		slog.String("func", "FollowUser"),
	)

	userID, err := uuid.Parse(ctx.Param("userId"))
	if err != nil {
		log.Warn("Error to parse UUID", slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid UUID", "The ID provided is not a valid UUID.")
	}

	if err := f.followerService.FollowUser(ctx.Request().Context(), userID); err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrFollowerNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The follower does not exist.")
		}

		if err == domain.ErrFollowerAlreadyExists {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, nil, "Conflict", "The follower already exists.")
		}

		if err == domain.ErrUserCannotFollowItself {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, nil, "Conflict", "The user cannot follow itself.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (f *followerHandler) UnfollowUser(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "follower"),
		slog.String("func", "UnfollowUser"),
	)

	followerID, err := uuid.Parse(ctx.Param("userId"))
	if err != nil {
		log.Warn("Error to parse UUID", slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid UUID", "The ID provided is not a valid UUID.")
	}

	if err := f.followerService.UnfollowUser(ctx.Request().Context(), followerID); err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrFollowerNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The follower does not exist.")
		}

		if err == domain.ErrFollowingNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The following does not exist.")
		}

		if err == domain.ErrUserCannotUnfollowItself {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, nil, "Conflict", "The user cannot unfollow itself.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (f *followerHandler) GetFollowers(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "follower"),
		slog.String("func", "GetFollowers"),
	)

	response, err := f.followerService.GetFollowers(ctx.Request().Context())
	if err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrUserNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The user does not exist.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (f *followerHandler) GetFollowings(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "follower"),
		slog.String("func", "GetFollowings"),
	)

	response, err := f.followerService.GetFollowings(ctx.Request().Context())
	if err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrUserNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The user does not exist.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}
