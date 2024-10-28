package main

import (
	"context"
	"log"
	"time"

	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/database"
	"github.com/G-Villarinho/social-network/domain"
)

func main() {
	config.LoadEnvironments()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewMysqlConnection(ctx)
	if err != nil {
		log.Fatal("error to connect to mysql: ", err)
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Follower{},
		&domain.Post{},
		&domain.Like{},
	); err != nil {
		log.Fatal("error to migrate: ", err)
	}

	log.Println("Migration executed successfully")
}
