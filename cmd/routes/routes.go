package routes

import (
	"github.com/gofiber/fiber/v2"
	"omsoft.com/auth/cmd/controllers"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/", status)

	registerAuth(v1)

	// app.Get("/api/bookmark", bookmark.GetAllBookmarks)
	// app.Post("/api/bookmark", bookmark.SaveBookmark)
}

func registerAuth(api fiber.Router) {
	auth := api.Group("/auth")
	auth.Post("/login", controllers.Login)
	// auth.Use(middleware.Protected())
	// auth.Get("/logout", Controllers.Logout(db))
}

func status(c *fiber.Ctx) error {
	return c.SendString("Server is running!")
}


