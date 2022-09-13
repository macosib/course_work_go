package main

import (
	"Attestation_work/pkg/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
)

func main() {
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

	server := &http.Server{Addr: "localhost:8000", Handler: router}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("не удалось запустить сервер: ", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	handler.Storage.WriteToCsv()
	server.Shutdown(ctx)
}
