package routes

import (
	"github.com/gofiber/fiber/v2"
	"omsoft.com/auth/cmd/controllers"
	"omsoft.com/auth/cmd/database"
)

func SetupRoutes(app *fiber.App, db *database.Database) {

	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/", status)

	registerAuth(v1, db)
	registerUser(v1, db)

	// app.Get("/api/bookmark", bookmark.GetAllBookmarks)
	// app.Post("/api/bookmark", bookmark.SaveBookmark)
}

func registerUser(api fiber.Router, db *database.Database) {
	user := api.Group("/user")
	user.Get("/profile1", controllers.Profile1)
	user.Use(controllers.AuthorizationRequired())
	user.Get("/profile2", controllers.Profile2)
	// auth.Use(middleware.Protected())
	// auth.Get("/logout", Controllers.Logout(db))
}

func registerAuth(api fiber.Router, db *database.Database) {
	auth := api.Group("/auth")
	auth.Post("/login", controllers.Login)

	// auth.Get("/profile1", controllers.Profile1)
	// auth.Use(controllers.AuthorizationRequired())
	// auth.Get("/profile2", controllers.Profile2)
	// auth.Use(middleware.Protected())
	// auth.Get("/logout", Controllers.Logout(db))
}

func status(c *fiber.Ctx) error {
	return c.SendString("Server is running!")
}
