package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/zigbee-s/go-todo-mongo-api/database"
	"github.com/zigbee-s/go-todo-mongo-api/environment"
	"github.com/zigbee-s/go-todo-mongo-api/router"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
	fmt.Println("Finished")
}

func run() error {
	// init env
	err := environment.LoadEnv()
	if err != nil {
		return err
	}

	// init db
	err = database.InitDB()
	if err != nil {
		return err
	}

	// defer closing db
	defer database.CloseDB()

	// create app
	app := fiber.New()

	// add basic middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// add routes
	router.AddTodoGroup(app)

	// start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))

	return nil
}
