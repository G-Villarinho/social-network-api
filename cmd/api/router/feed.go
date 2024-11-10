package router

import (
	"log"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/labstack/echo/v4"
)

func setupFeedRoutes(e *echo.Echo, di *internal.Di) {
	feedHandler, err := internal.Invoke[domain.FeedHandler](di)
	if err != nil {
		log.Fatal("error to create user handler: ", err)
	}

	group := e.Group("/v1/feed")

	group.GET("", feedHandler.GetFeed)
}
