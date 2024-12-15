package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/khanirfan96/To-do-Fullstack.git/router"
)

func main() {
	fmt.Println("FullStack TODO Application")

	r := router.Router()
	fmt.Println("Server is getting Started.....")
	log.Fatal(http.ListenAndServe(":8000", r))
	fmt.Println("Server is started at port 8000.....")
}
