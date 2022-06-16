package controllers

import (
	"github.com/gofiber/fiber/v2"
	"omsoft.com/auth/cmd/models"
)

type (
	MsgLogin models.Login
)

func Login(c *fiber.Ctx) error {
	var l MsgLogin
	err := c.BodyParser(&l)
	if err != nil {
		// return failOnError(c, err, "cannot parse json", fiber.StatusBadRequest)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"message": err,
				"data":    nil,
			})
		}
	}
	println(l.Username)
	println(l.Password)
	return nil
}
