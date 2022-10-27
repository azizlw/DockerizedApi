package router

import (
	"github.com/azizlw/FinalProject/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// register
	router.HandleFunc("/api/post/register", controller.Register).Methods("POST")

	// login & generating token and returning it in response
	router.HandleFunc("/api/post/login", controller.Login).Methods("POST")

	// storing details in cart for a particular user
	router.HandleFunc("/api/post/user/{id}", controller.UserCart).Methods("POST")

	// fetching user cart details
	router.HandleFunc("/api/get/getcart/{username}", controller.GetCart).Methods("GET")

	router.HandleFunc("/api/get/home/items", controller.GetAllItems).Methods("GET")
	router.HandleFunc("/api/post/item", controller.InsertOneItem).Methods("POST")
	router.HandleFunc("/api/put/item/{id}", controller.UpdateOneItem).Methods("PUT")
	router.HandleFunc("/api/delete/item/{id}", controller.DeleteOneItem).Methods("DELETE")
	router.HandleFunc("/api/delete/deleteallitems", controller.DeleteAllItems).Methods("DELETE")

	return router
}
