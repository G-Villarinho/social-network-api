package middleware

import (
	"context"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/labstack/echo/v4"
)

func ClientInfo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Echo().IPExtractor == nil {
			c.Echo().IPExtractor = echo.ExtractIPFromXFFHeader(
				echo.TrustLoopback(false),   // e.g. ipv4 start with 127.
				echo.TrustLinkLocal(false),  // e.g. ipv4 start with 169.254
				echo.TrustPrivateNet(false), // e.g. ipv4 start with 10. or 192.168
			)
		}

		userAgent := c.Request().Header.Get("User-Agent")
		clientIP := c.RealIP()

		ctx := context.WithValue(c.Request().Context(), domain.UserAgentKey, userAgent)
		ctx = context.WithValue(ctx, domain.ClientIPKey, clientIP)

		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
