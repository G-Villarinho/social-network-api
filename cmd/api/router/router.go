package router

import (
	"github.com/G-Villarinho/social-network/internal"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, di *internal.Di) {
	setupUserRoutes(e, di)
	setupFollowerRoutes(e, di)
	setupPostRoutes(e, di)
}
