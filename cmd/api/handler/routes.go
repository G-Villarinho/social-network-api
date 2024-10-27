package handler

import (
	"log"

	"github.com/G-Villarinho/social-network/pkg"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, di *pkg.Di) {
	setupUserRoutes(e, di)
}

func setupUserRoutes(e *echo.Echo, di *pkg.Di) {
	userHandler, err := NewUserHandler(di)
	if err != nil {
		log.Fatal("error to create user handler: ", err)
	}

	group := e.Group("/v1/users")

	group.POST("", userHandler.CreateUser)
}
