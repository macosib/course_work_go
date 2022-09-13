package main

import (
	"Attestation_work/pkg/handlers"
	"net/http"
	"github.com/go-chi/chi"
)

func main () {
	handler := handlers.GetHandler()
	router := chi.NewRouter()
	router.Route("/api/v1/city", func(router chi.Router) {
		router.Get("/", handler.GetInfoCityView)
		router.Post("/", handler.AddCityView)
		router.Route("/{Id}", func(router chi.Router) {
			router.Get("/", handler.CityView)    
			router.Delete("/", handler.CityView) 
			router.Patch("/", handler.CityView) 
	})
	})
	
	http.ListenAndServe("localhost:8000", router)
}
