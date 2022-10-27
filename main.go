package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/azizlw/FinalProject/router"
)

func main() {
	fmt.Println("Inventory API with MongoDB")
	r := router.Router()
	fmt.Println("Server is getting started...")
	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listening at port 4000...")
}
