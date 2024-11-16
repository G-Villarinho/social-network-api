package router

import (
	"log"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/G-Villarinho/social-network/middleware"

	"github.com/labstack/echo/v4"
)

func setupUserRoutes(e *echo.Echo, di *internal.Di) {
	userHandler, err := internal.Invoke[domain.UserHandler](di)
	if err != nil {
		log.Fatal("error to create user handler: ", err)
	}

	group := e.Group("/v1/users")

	group.POST("", userHandler.CreateUser)
	group.POST("/sign-in", userHandler.SignIn, middleware.ClientInfo)
	group.POST("/sign-out", userHandler.SignOut, middleware.EnsureAuthenticated(di))
	group.GET("/me", userHandler.GetUser, middleware.EnsureAuthenticated(di))
	group.PUT("", userHandler.UpdateUser, middleware.EnsureAuthenticated(di))
	group.DELETE("", userHandler.DeleteUser, middleware.EnsureAuthenticated(di))
	group.POST("/check-username", userHandler.CheckUsername)
	group.POST("/check-password-strong", userHandler.CheckPasswordStrong)
}
