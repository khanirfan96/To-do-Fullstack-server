package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/khanirfan96/To-do-Fullstack-server/middleware"
)

func Router() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "*",
	}))

	api := app.Group("/api")
	recipeapi := app.Group("/recipe")

	api.Get("/gettodo", middleware.GetTodo)
	api.Post("/posttodo", middleware.CreateTodo)
	api.Put("/puttodo/:id", middleware.UpdateTodo)
	api.Put("/undotodo/:id", middleware.UndoTodo)
	api.Delete("/deleteonetodo/:id", middleware.DeleteOneTodo)
	api.Delete("/deletetodo", middleware.DeleteAllTodo)

	// *********************** recipe API ******************************

	recipeapi.Get("/getrecipe", middleware.GetRecipe)
	recipeapi.Post("/postrecipe", middleware.CreateRecipe)
	recipeapi.Put("/putrecipe/:id", middleware.UpdateRecipe)
	recipeapi.Put("/putingredients/:id", middleware.UpdateIngredeints)
	recipeapi.Delete("/deleterecipe/:id", middleware.DeleteOneRecipe)
	recipeapi.Delete("/deleterecipe", middleware.DeleteAllRecipe)

	return app
}
