package util

import "github.com/gofiber/fiber/v2"

func SuccessResponse[T any](c *fiber.Ctx, msg string, status int, v T) error {
	c.Status(status).JSON(fiber.Map{
		"msg":   msg,
		"error": nil,
		"data":  v,
	})
	return nil
}

func FailOnError[T any](c *fiber.Ctx, err error, msg string, status int, v T) error {
	if err != nil {
		c.Status(status).JSON(fiber.Map{
			"msg":   msg,
			"error": err.Error(),
			"data":  v,
		})
	}
	return nil
}
