package main

import (
	"fmt"
	"log"

	"github.com/khanirfan96/To-do-Fullstack-server/router"
)

func main() {
	fmt.Println("FullStack TODO Application")
	r := router.Router()

	fmt.Println("Server is getting Started.....")
	log.Fatal(r.Listen(":8000"))
	fmt.Println("Server is started at port 8000.....")
}
