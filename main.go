package main

import (
	"context"

	"github.com/conflux-tech/fiber-rest-boilerplate/configs"
	"github.com/conflux-tech/fiber-rest-boilerplate/database"
	"github.com/conflux-tech/fiber-rest-boilerplate/routes"
	"github.com/conflux-tech/fiber-rest-boilerplate/users"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	ctx := context.Background()
	conf, _ := configs.Load()

	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(compress.New())
	app.Use(requestid.New())

	if conf.Database.Enabled {
		database.Connect(&conf.Database)
		defer database.Close()
	}

	redis := setupRedis(ctx)

	userPGRepo := users.NewPGRepo(database.Instance())
	userUsecase := users.NewUseCase(userPGRepo, *redis)

	routes.RegisterRoot(app.Group("/"))
	routes.RegisterUsers(ctx, app.Group("/users"), *userUsecase)

	app.Listen(conf.Application.Port)
}

func setupRedis(ctx context.Context) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: configs.Get().Redis.Host,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	return rdb
}
