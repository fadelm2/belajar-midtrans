package main

import (
	"github.com/fadelm2/belajar_midtrans/config"
	"github.com/fadelm2/belajar_midtrans/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.ConnectDB()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://192.168.1.70:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	app.Post("/checkout", handlers.Checkout)
	app.Post("/webhook/midtrans", handlers.MidtransWebhookHandler)

	app.Listen(":8080")
}
