package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ServiceConfig struct {
	Port       string
	MongoDBURL string
}

func main() {
	fmt.Println("todo-list service started...")

	config := ServiceConfig{
		Port:       ":8080",
		MongoDBURL: "mongodb://localhost:27017",
	}
	repository := NewRepository(config.MongoDBURL)
	service := NewService(repository)
	api := NewAPI(service)
	app := ServiceSetup(api)

	app.Listen(config.Port)
}

func ServiceSetup(api *Api) *fiber.App {
	app := fiber.New()
	app.Post("/todo", api.PostTodoApi)
	app.Get("/todo", api.GetTodoListApi)
	app.Get("/todo/:id", api.GetTodoApi)
	app.Put("/todo/:id", api.PutTodoApi)
	app.Put("/sort", api.PutSortApi)
	app.Delete("/todo/:id", api.DeleteTodoApi)

	return app
}
