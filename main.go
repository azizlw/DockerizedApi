package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/azizlw/FinalProject/router"
)

func main() {
	fmt.Println("Inventory API with MongoDB")
	r := router.Router()
	fmt.Println("Server is getting started...")
	// log.Fatal(http.ListenAndServe(":4000", r))
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), r))
	fmt.Println("Listening at port 4000...")
}
