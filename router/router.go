package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/khanirfan96/To-do-Fullstack.git/middleware"
	"github.com/rs/cors"
)

func Router() http.Handler {
	route := mux.NewRouter()

	router := route.PathPrefix("/api").Subrouter()

	router.HandleFunc("/gettodo", middleware.GetTodo).Methods("GET")
	router.HandleFunc("/posttodo", middleware.CreateTodo).Methods("POST")
	router.HandleFunc("/puttodo/{id}", middleware.UpdateTodo).Methods("PUT")
	router.HandleFunc("/undotodo/{id}", middleware.UndoTodo).Methods("PUT")
	router.HandleFunc("/deleteonetodo/{id}", middleware.DeleteOneTodo).Methods("DELETE")
	router.HandleFunc("/deletetodo", middleware.DeleteAllTodo).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"}, // Allow all headers
		Debug:          true,          // Enable debugging for testing, remove in production
	})

	// Wrap the router with the CORS handler
	corsHandle := c.Handler(router)

	return corsHandle
}
