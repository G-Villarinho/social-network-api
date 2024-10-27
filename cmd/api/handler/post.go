package handler

import (
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
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
