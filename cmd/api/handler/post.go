package handler

import (
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type postHandler struct {
	di          *pkg.Di
	postService domain.PostService
}

func NewPostHandler(di *pkg.Di) (domain.PostHandler, error) {
	postService, err := pkg.Invoke[domain.PostService](di)
	if err != nil {
		return nil, err
	}

	return &postHandler{
		di:          di,
		postService: postService,
	}, nil
}

func (p *postHandler) CreatePost(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "post"),
		slog.String("func", "CreatePost"),
	)

	var payload domain.PostPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Error to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	if err := p.postService.CreatePost(ctx.Request().Context(), payload); err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusCreated)
}

func (p *postHandler) GetPosts(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "post"),
		slog.String("func", "GetPosts"),
	)

	response, err := p.postService.GetPosts(ctx.Request().Context())
	if err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrPostNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The post does not exist.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (p *postHandler) GetPostById(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "post"),
		slog.String("func", "GetPostById"),
	)

	ID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Warn("Error to parse UUID", slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Bad Request", "Invalid ID.")
	}

	response, err := p.postService.GetPostById(ctx.Request().Context(), ID)
	if err != nil {
		log.Error(err.Error())

		if err == domain.ErrPostNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The post does not exist.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (p *postHandler) UpdatePost(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "post"),
		slog.String("func", "UpdatePost"),
	)

	ID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Warn("Error to parse UUID", slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Bad Request", "Invalid ID.")
	}

	var payload domain.PostUpdatePayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Error to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	if err := p.postService.UpdatePost(ctx.Request().Context(), ID, payload); err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrPostNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The post does not exist.")
		}

		if err == domain.ErrPostNotBelongToUser {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, nil, "Conflict", "The post does not belong to the user.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusOK)
}

func (p *postHandler) DeletePost(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "post"),
		slog.String("func", "DeletePost"),
	)

	ID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Warn("Error to parse UUID", slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Bad Request", "Invalid ID.")
	}

	if err := p.postService.DeletePost(ctx.Request().Context(), ID); err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrPostNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The post does not exist.")
		}

		if err == domain.ErrPostNotBelongToUser {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, nil, "Conflict", "The post does not belong to the user.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (p *postHandler) GetByUserID(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "post"),
		slog.String("func", "GetByUserID"),
	)

	userID, err := uuid.Parse(ctx.Param("userId"))
	if err != nil {
		log.Warn("Error to parse UUID", slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Bad Request", "Invalid ID.")
	}

	response, err := p.postService.GetByUserID(ctx.Request().Context(), userID)
	if err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrPostNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The post does not exist.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)

}

func (p *postHandler) LikePost(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "post"),
		slog.String("func", "LikePost"),
	)

	ID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Warn("Error to parse UUID", slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Bad Request", "Invalid ID.")
	}

	if err := p.postService.LikePost(ctx.Request().Context(), ID); err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrPostNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The post does not exist.")
		}

		if err == domain.ErrPostAlreadyLiked {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, nil, "Conflict", "The post is already liked.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusOK)
}

func (p *postHandler) UnlikePost(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "post"),
		slog.String("func", "UnLikePost"),
	)

	ID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Warn("Error to parse UUID", slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Bad Request", "Invalid ID.")
	}

	if err := p.postService.UnlikePost(ctx.Request().Context(), ID); err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrPostNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The post does not exist.")
		}

		if err == domain.ErrPostNotLiked {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, nil, "Conflict", "The post is not liked.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusOK)
}
