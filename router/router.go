package router

import (
	"github.com/gorilla/mux"
	"github.com/khanirfan96/To-do-Fullstack.git/middleware"
)

func Router() *mux.Router {
	route := mux.NewRouter()

	router := route.PathPrefix("/api").Subrouter()

	router.HandleFunc("/gettodo", middleware.GetTodo).Methods("GET", "OPTIONS")
	router.HandleFunc("/posttodo", middleware.CreateTodo).Methods("POST", "OPTIONS")
	router.HandleFunc("/puttodo/{id}", middleware.UpdateTodo).Methods("PUT", "OPTIONS")
	router.HandleFunc("/undotodo/{id}", middleware.UndoTodo).Methods("PUT", "OPTIONS")
	router.HandleFunc("/deleteonetodo/{id}", middleware.DeleteOneTodo).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/deletetodo", middleware.DeleteAllTodo).Methods("DELETE", "OPTIONS")

	return router
}
