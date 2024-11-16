package middleware

import (
	"context"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/labstack/echo/v4"
)

func ClientInfo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userAgent := c.Request().Header.Get("User-Agent")
		clientIP := c.RealIP()

		ctx := context.WithValue(c.Request().Context(), domain.UserAgentKey, userAgent)
		ctx = context.WithValue(ctx, domain.ClientIPKey, clientIP)

		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
