package handler

import (
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	di          *pkg.Di
	userService domain.UserService
}

func NewUserHandler(di *pkg.Di) (domain.UserHandler, error) {
	userService, err := pkg.Invoke[domain.UserService](di)
	if err != nil {
		return nil, err
	}

	return &userHandler{
		di:          di,
		userService: userService,
	}, nil
}

func (u *userHandler) CreateUser(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "CreateUser"),
	)

	var payload domain.UserPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Error to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	token, err := u.userService.CreateUser(ctx.Request().Context(), payload)
	if err != nil {
		log.Error(err.Error())
		if err == domain.ErrEmailAlreadyRegister {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, nil, "conflict", "The email already registered. Please try again with a different email.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	cookie := &http.Cookie{
		Name:     "x.Token",
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusCreated)
}

func (u *userHandler) SignIn(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "SignIn"),
	)

	var payload domain.SignInPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("error to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	token, err := u.userService.SignIn(ctx.Request().Context(), payload)
	if err != nil {
		if err == domain.ErrUserNotFound || err == domain.ErrInvalidPassword {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusUnauthorized, nil, "Unauthorized", "Invalid email or password. Please check your credentials and try again.")
		}

		log.Error(err.Error())
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	cookie := &http.Cookie{
		Name:     "x.Token",
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}

func (u *userHandler) SignOut(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "SignOut"),
	)

	if err := u.userService.SignOut(ctx.Request().Context()); err != nil {
		log.Error(err.Error())

		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	cookie := &http.Cookie{
		Name:     "x.Token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}

func (u *userHandler) GetUser(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "GetUser"),
	)

	response, err := u.userService.GetUser(ctx.Request().Context())
	if err != nil {
		log.Error(err.Error())
		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrUserNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "not_found", "User not found.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (u *userHandler) UpdateUser(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "UpdateUser"),
	)

	var payload domain.UserUpdatePayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("error to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	if err := u.userService.UpdateUser(ctx.Request().Context(), payload); err != nil {
		log.Error(err.Error())
		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		if err == domain.ErrUserNotFound {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "not_found", "User not found.")
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusOK)
}

func (u *userHandler) DeleteUser(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "DeleteUser"),
	)

	if err := u.userService.DeleteUser(ctx.Request().Context()); err != nil {
		log.Error(err.Error())
		if err == domain.ErrSessionNotFound {
			return domain.AccessDeniedAPIErrorResponse(ctx)
		}

		return domain.InternalServerAPIErrorResponse(ctx)
	}

	cookie := &http.Cookie{
		Name:     "x.Token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}
