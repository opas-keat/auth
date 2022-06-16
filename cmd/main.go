package main

import (
	"os"
	"os/signal"

	configuration "omsoft.com/auth/cmd/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"omsoft.com/auth/cmd/database"
	"omsoft.com/auth/cmd/routes"
)

type App struct {
	*fiber.App
	DB *database.Database
}

func main() {
	// // JWT Middleware
	// // app.Use(jwtware.New(jwtware.Config{
	// // 	SigningKey: []byte("omsoft"),
	// // }))

	// // db, err := database.New(&database.DatabaseConfig{
	// // 	Driver:   config.GetString("DB_DRIVER"),
	// // 	Host:     config.GetString("DB_HOST"),
	// // 	Username: config.GetString("DB_USERNAME"),
	// // 	Password: config.GetString("DB_PASSWORD"),
	// // 	Port:     config.GetInt("DB_PORT"),
	// // 	Database: config.GetString("DB_DATABASE"),
	// // })

	config := configuration.New()

	app := App{
		App: fiber.New(*config.GetFiberConfig()),
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept,Authorization",
	}))
	// app.Use(csrf.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(logger.New(logger.Config{
		Format:     "${cyan}[${time}] ${white}${pid} ${red}${status} ${blue}[${method}] ${white}${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}))

	routes.SetupRoutes(app.App)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		app.exit()
	}()
	// Start listening on the specified address
	err := app.Listen(config.GetString("APP_ADDR"))
	if err != nil {
		app.exit()
	}
}

func (app *App) exit() {
	_ = app.Shutdown()
}
