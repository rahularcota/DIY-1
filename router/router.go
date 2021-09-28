package router

import (
	"github.com/gorilla/mux"
	"github.com/rahularcota/DIY1/handlers"
)

func InitializeRouter() *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc("/product", handlers.CreateProduct).Methods("POST")
	router.HandleFunc("/product/{id:[0-9]+}", handlers.GetProduct).Methods("GET")
	router.HandleFunc("/product/{id:[0-9]+}", handlers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/product/{id:[0-9]+}", handlers.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/store/{id}/products", handlers.GetStoreProducts).Methods("GET")
	router.HandleFunc("/store/{id}", handlers.AddStoreProducts).Methods("POST")
	router.HandleFunc("/store", handlers.CreateStore).Methods("POST")

	return router
}
