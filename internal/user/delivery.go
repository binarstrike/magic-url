package user

import "github.com/gofiber/fiber/v2"

type Controller interface {
	Register(*fiber.Ctx) error
	Login(*fiber.Ctx) error
	Logout(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
	Me(*fiber.Ctx) error
}
