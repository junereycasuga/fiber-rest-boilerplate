package routes

import (
	"context"

	"github.com/conflux-tech/fiber-rest-boilerplate/handlers"
	"github.com/conflux-tech/fiber-rest-boilerplate/users"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoot registers all available routes for root
func RegisterRoot(route fiber.Router) {
	route.Get("/", handlers.GetStatus)
}

// RegisterUsers registers all available routes for /user
func RegisterUsers(ctx context.Context, route fiber.Router, usecase users.UseCase) {
	route.Get("/", handlers.ListUsers(ctx, usecase))
	route.Get("/:id", handlers.GetUser(ctx, usecase))
	route.Post("/", handlers.CreateUser(ctx, usecase))
	route.Patch("/:id", handlers.UpdateUser(ctx, usecase))
	route.Delete("/:id", handlers.DeleteUser(ctx, usecase))
}
