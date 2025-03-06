package middleware

import (
	"log"

	"github.com/khanirfan96/To-do-Fullstack-server/database"
)

func init() {
	if err := database.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
}
