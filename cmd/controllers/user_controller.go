package controllers

import (
	"github.com/gofiber/fiber/v2"
	"omsoft.com/auth/cmd/models"
	"omsoft.com/auth/cmd/util"
)

func Profile1(c *fiber.Ctx) error {
	println("Profile1")
	return util.SuccessResponse(c, "", fiber.StatusOK, models.User{})
}

func Profile2(c *fiber.Ctx) error {
	println("Profile2")
	return util.SuccessResponse(c, "", fiber.StatusOK, models.User{})
}
