package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/conflux-tech/fiber-rest-boilerplate/users"
	"github.com/gofiber/fiber/v2"
)

// ListUsers returns all users from db
func ListUsers(ctx context.Context, usecase users.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Params("page"))
		limit, _ := strconv.Atoi(c.Params("limit"))
		users, err := usecase.List(ctx, page, limit)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"message": "An error occured while retrieiving data",
					"details": err,
				},
			})
		}
		return c.Status(http.StatusOK).JSON(users)
	}
}

// GetUser returns information about a user
func GetUser(ctx context.Context, usecase users.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		user, err := usecase.Get(ctx, id)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"message": "An error occured while retrieiving data",
					"details": err,
				},
			})
		}
		if user.ID == 0 {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"message": "User is invalid",
				},
			})
		}
		return c.Status(http.StatusOK).JSON(&user)
	}
}

// CreateUser creates new user record
func CreateUser(ctx context.Context, usecase users.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(users.User)
		if err := c.BodyParser(user); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(fiber.Map{
				"error": fiber.Map{
					"message": "An error occurred while parsing the request",
					"details": err,
				},
			})
		}
		u, err := usecase.Create(ctx, user)
		if err != nil {
			return fiber.NewError(000, "An error occured while creating record")
		}
		return c.Status(http.StatusCreated).JSON(u)
	}
}

// UpdateUser updates a user record
func UpdateUser(ctx context.Context, usecase users.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		editUser := new(users.User)

		if err := c.BodyParser(editUser); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(fiber.Map{
				"error": fiber.Map{
					"message": "An error occurred while parsing the request",
					"details": err,
				},
			})
		}

		u, err := usecase.Update(ctx, id, editUser)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"message": "An error occurred while updating the data",
					"details": err,
				},
			})
		}
		return c.JSON(u)
	}
}

// DeleteUser deletes a use
func DeleteUser(ctx context.Context, usecase users.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		res, err := usecase.Delete(ctx, id)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"message": "An error occured while deleting the data",
					"details": err,
				},
			})
		}
		return c.JSON(fiber.Map{
			"id":      id,
			"deleted": res,
		})
	}
}
