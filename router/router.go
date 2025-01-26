package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/khanirfan96/To-do-Fullstack.git/middleware"
)

func Router() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "*",
	}))

	api := app.Group("/api")

	api.Get("/gettodo", middleware.GetTodo)
	api.Post("/posttodo", middleware.CreateTodo)
	api.Put("/puttodo/:id", middleware.UpdateTodo)
	api.Put("/undotodo/:id", middleware.UndoTodo)
	api.Delete("/deleteonetodo/:id", middleware.DeleteOneTodo)
	api.Delete("/deletetodo", middleware.DeleteAllTodo)

	return app
}
