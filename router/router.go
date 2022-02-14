package router

import (
	"CRUD-app/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

    router := mux.NewRouter()

    router.HandleFunc("/api/harga/{id}", middleware.GetHarga).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/harga", middleware.GetAllHarga).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/newharga", middleware.CreateHarga).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/harga/{id}", middleware.UpdateHarga).Methods("PUT", "OPTIONS")
    router.HandleFunc("/api/deleteharga/{id}", middleware.DeleteHarga).Methods("DELETE", "OPTIONS")

    return router
}