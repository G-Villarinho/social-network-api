package router

import (
	"log"
	"net/http"
	"time"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/G-Villarinho/social-network/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"golang.org/x/time/rate"
)

func setupPostRoutes(e *echo.Echo, di *internal.Di) {
	postHandler, err := internal.Invoke[domain.PostHandler](di)
	if err != nil {
		log.Fatal("error to create user handler: ", err)
	}

	config := echomiddleware.RateLimiterConfig{
		Skipper: echomiddleware.DefaultSkipper,
		Store: echomiddleware.NewRateLimiterMemoryStoreWithConfig(
			echomiddleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	group := e.Group("/v1/posts", middleware.EnsureAuthenticated(di))

	group.POST("", postHandler.CreatePost)
	group.GET("/feed", postHandler.GetPosts)
	group.GET("/:id", postHandler.GetPostById)
	group.PUT("/:id", postHandler.UpdatePost)
	group.DELETE("/:id", postHandler.DeletePost)
	group.GET("/user/:userId", postHandler.GetByUserID)
	group.POST("/:id/like", postHandler.LikePost, echomiddleware.RateLimiterWithConfig(config))
	group.DELETE("/:id/like", postHandler.UnlikePost, echomiddleware.RateLimiterWithConfig(config))
}
