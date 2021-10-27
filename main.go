package main

import (
	"fmt"
	"log"
	"os"

	"linktory/database"
	"linktory/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func startEnv() {
    // loads values from .env into the system
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
	log.Print("env found")
}


func main() {
	// startEnv()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	fmt.Println("Environment variables successfully loaded. Starting application...")
	database.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	// PORT := os.Getenv("PORT") 
	port, ok := os.LookupEnv("PORT")

	// if there is no PORT in environment variables default to port 8000
	if !ok {
		port = "8000"
	}

	fmt.Println("Application started...")
	app.Listen(fmt.Sprintf(":%s", port))
}