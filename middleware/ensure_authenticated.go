package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/labstack/echo/v4"
)

func EnsureAuthenticated(di *pkg.Di) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			sessionService, err := pkg.Invoke[domain.SessionService](di)
			if err != nil {
				slog.Error(err.Error())
				return domain.InternalServerAPIErrorResponse(ctx)
			}

			cookie, err := ctx.Cookie("x.Token")
			if err != nil || cookie == nil || cookie.Value == "" {
				clearAuthCookie(ctx)
				return domain.AccessDeniedAPIErrorResponse(ctx)
			}

			session, err := sessionService.GetSessionByToken(ctx.Request().Context(), cookie.Value)
			if err != nil {
				if err == domain.ErrTokenInvalid || err == domain.ErrSessionMismatch || err == domain.ErrSessionNotFound {
					clearAuthCookie(ctx)
					return domain.AccessDeniedAPIErrorResponse(ctx)
				}
				slog.Error(err.Error())
				return domain.InternalServerAPIErrorResponse(ctx)
			}

			newCtx := context.WithValue(ctx.Request().Context(), domain.SessionKey, session)
			ctx.SetRequest(ctx.Request().WithContext(newCtx))

			return next(ctx)
		}
	}
}

func clearAuthCookie(ctx echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "x.Token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = -1

	ctx.SetCookie(cookie)
}
