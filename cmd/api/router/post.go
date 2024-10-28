package router

import (
	"log"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/middleware"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/labstack/echo/v4"
)

func setupPostRoutes(e *echo.Echo, di *pkg.Di) {
	postHandler, err := pkg.Invoke[domain.PostHandler](di)
	if err != nil {
		log.Fatal("error to create user handler: ", err)
	}

	group := e.Group("/v1/posts", middleware.EnsureAuthenticated(di))

	group.POST("", postHandler.CreatePost)
	group.GET("", postHandler.GetPosts)
	group.GET("/:id", postHandler.GetPostById)
	group.PUT("/:id", postHandler.UpdatePost)
	group.DELETE("/:id", postHandler.DeletePost)
	group.GET("/user/:userId", postHandler.GetByUserID)
	group.POST("/:id/like", postHandler.LikePost)
	group.DELETE("/:id/like", postHandler.UnlikePost)
}
