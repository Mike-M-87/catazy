package main

import (
	"cataz/cataz"
	"cataz/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

var defaultPort = "8082"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic().Err(err)
	}
	utils.InitialiseDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	app := fiber.New(fiber.Config{
		Prefork: false,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Get("/movies", cataz.FetchMovies)
	app.Post("/watch", cataz.Watch)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(418).JSON(&fiber.Map{
			"Message": "🍏 Rada bro, umepotea kwani",
		}) // => 418 "I am a tepot"
	})

	initCron()

	log.Fatal().Err(app.Listen(":" + port))
}

func initCron() {
	c := cron.New()
	_, err := c.AddFunc("@hourly", cataz.StoreMovies)
	if err != nil {
		return
	}
	println("🕓 Cron start")
	c.Start()
}
