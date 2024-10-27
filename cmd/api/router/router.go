package router

import (
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, di *pkg.Di) {
	setupUserRoutes(e, di)
}
