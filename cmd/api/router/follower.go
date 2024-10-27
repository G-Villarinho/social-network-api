package router

import (
	"log"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/middleware"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/labstack/echo/v4"
)

func setupFollowerRoutes(e *echo.Echo, di *pkg.Di) {
	userHandler, err := pkg.Invoke[domain.FollowerHandler](di)
	if err != nil {
		log.Fatal("error to create follower handler: ", err)
	}

	group := e.Group("/v1/followers", middleware.EnsureAuthenticated(di))

	group.POST("/:userId", userHandler.FollowUser)
	group.DELETE("/:userId", userHandler.UnfollowUser)
	group.GET("", userHandler.GetFollowers)
	group.GET("/fowllings", userHandler.GetFollowings)

}
