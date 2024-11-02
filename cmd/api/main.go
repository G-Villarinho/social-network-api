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
	"github.com/G-Villarinho/social-network/pkg" // ajuste conforme necessário
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
	di := pkg.NewDi()

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

	pkg.Provide(di, func(d *pkg.Di) (*gorm.DB, error) {
		return db, nil
	})

	pkg.Provide(di, func(d *pkg.Di) (*redis.Client, error) {
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

	pkg.Provide(di, func(d *pkg.Di) (client.RabbitMQClient, error) {
		return rabbitMQClient, nil
	})

	pkg.Provide(di, handler.NewFollowerHandler)
	pkg.Provide(di, handler.NewPostHandler)
	pkg.Provide(di, handler.NewUserHandler)

	pkg.Provide(di, service.NewContextService)
	pkg.Provide(di, service.NewFollowerService)
	pkg.Provide(di, service.NewPostService)
	pkg.Provide(di, service.NewQueueService)
	pkg.Provide(di, service.NewSessionService)
	pkg.Provide(di, service.NewUserService)

	pkg.Provide(di, repository.NewFollowerRepository)
	pkg.Provide(di, repository.NewMemoryCacheRepository)
	pkg.Provide(di, repository.NewPostRepository)
	pkg.Provide(di, repository.NewSessionRepository)
	pkg.Provide(di, repository.NewUserRepository)

	router.SetupRoutes(e, di)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Env.APIPort)))
}
