package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/G-Villarinho/social-network/client"
	"github.com/G-Villarinho/social-network/cmd/api/handler"
	"github.com/G-Villarinho/social-network/cmd/api/router"
	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/database"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/G-Villarinho/social-network/repository"
	"github.com/G-Villarinho/social-network/service"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func main() {
	config.ConfigureLogger()
	config.LoadEnvironments()

	e := echo.New()
	di := internal.NewDi()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","method":"${method}","uri":"${uri}","status":${status},"latency":"${latency_human}"}`,
		Output: os.Stdout,
	}))

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		fmt.Printf("\n")
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{config.Env.FrontURL},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewMysqlConnection(ctx)
	if err != nil {
		log.Fatal("error to connect to mysql: ", err)
	}

	redisClient, err := database.NewRedisConnection(ctx)
	if err != nil {
		log.Fatal("error to connect to redis: ", err)
	}

	internal.Provide(di, func(d *internal.Di) (*gorm.DB, error) {
		return db, nil
	})

	internal.Provide(di, func(d *internal.Di) (*redis.Client, error) {
		return redisClient, nil
	})

	rabbitMQClient, err := client.NewRabbitMQClient(di)
	if err != nil {
		log.Fatal("error initializing RabbitMQ client: ", err)
	}
	if err := rabbitMQClient.Connect(); err != nil {
		log.Fatal("error connecting to RabbitMQ: ", err)
	}
	defer func() {
		if err := rabbitMQClient.Disconnect(); err != nil {
			log.Println("error disconnecting from RabbitMQ:", err)
		}
	}()

	internal.Provide(di, func(d *internal.Di) (client.RabbitMQClient, error) {
		return rabbitMQClient, nil
	})

	internal.Provide(di, handler.NewFollowerHandler)
	internal.Provide(di, handler.NewPostHandler)
	internal.Provide(di, handler.NewUserHandler)

	internal.Provide(di, service.NewContextService)
	internal.Provide(di, service.NewFeedService)
	internal.Provide(di, service.NewFollowerService)
	internal.Provide(di, service.NewLikeService)
	internal.Provide(di, service.NewPostService)
	internal.Provide(di, service.NewQueueService)
	internal.Provide(di, service.NewSessionService)
	internal.Provide(di, service.NewUserService)

	internal.Provide(di, repository.NewFollowerRepository)
	internal.Provide(di, repository.NewLikeRepository)
	internal.Provide(di, repository.NewMemoryCacheRepository)
	internal.Provide(di, repository.NewPostRepository)
	internal.Provide(di, repository.NewSessionRepository)
	internal.Provide(di, repository.NewUserRepository)

	router.SetupRoutes(e, di)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Env.APIPort)))
}
