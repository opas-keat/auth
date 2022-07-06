package main

import (
	"fmt"
	"os"
	"os/signal"

	configuration "omsoft.com/auth/cmd/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"omsoft.com/auth/cmd/database"
	"omsoft.com/auth/cmd/routes"
)

type App struct {
	*fiber.App
	DB *database.Database
}

var (
	// initialVector = "4244467890218023942835864651981516513207895156132164643213169699"
	passphrase = "AAACCCDDDYYUURRS"
)

func main() {
	// // JWT Middleware
	// // app.Use(jwtware.New(jwtware.Config{
	// // 	SigningKey: []byte("omsoft"),
	// // }))

	config := configuration.New()

	db, err := database.New(&database.DatabaseConfig{
		Driver:   config.GetString("DB_DRIVER"),
		Host:     config.GetString("DB_HOST"),
		Username: config.GetString("DB_USERNAME"),
		Password: config.GetString("DB_PASSWORD"),
		Port:     config.GetInt("DB_PORT"),
		Database: config.GetString("DB_DATABASE"),
	})
	if err != nil {
		fmt.Println("failed to connect to database:", err.Error())
		if db == nil {
			fmt.Println("failed to connect to database: db variable is nil")
		}
	}

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
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	routes.SetupRoutes(app.App, db)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		app.exit()
	}()
	err = app.Listen(config.GetString("APP_ADDR"))
	if err != nil {
		app.exit()
	}
}

func (app *App) exit() {
	_ = app.Shutdown()
}

// package main

// import (
// 	"encoding/base64"

// 	"omsoft.com/auth/cmd/util"
// )

// func main() {
// 	srcData := "123456 Test ทดสอบ 001/2565 ๑๒๓๔๕๖ @<>-/*-"
// 	//Test encryption
// 	encData, err := util.ECBEncrypt([]byte(srcData))
// 	if err != nil {
// 		// fmt.Errorf(err.Error())
// 		return
// 	}
// 	encDataStr := base64.StdEncoding.EncodeToString(encData)
// 	println(encDataStr)
// 	// CPlpXAYFIjLgfhebVjUv9emDxXiy7qI39MQ/YHnFWSEaXUZYxriFCRntIVoJA9wMjpRe8QMvLHNTHJ4ooMtf+w==
// 	encData2, err := base64.StdEncoding.DecodeString(encDataStr)
// 	if err != nil {
// 		// fmt.Errorf(err.Error())
// 		return
// 	}
// 	//Test decryption
// 	decData, err := util.ECBDecrypt(encData2)
// 	if err != nil {
// 		// fmt.Errorf(err.Error())
// 		return
// 	}
// 	println(string(decData))
// }
