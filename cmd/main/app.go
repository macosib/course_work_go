package main

import (
	"Attestation_work/internal/city"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	storage := city.NewStore()
	router := httprouter.New()
	handler := city.NewHandler(storage)
	handler.Register(router)

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

	city.WriteToCsv(storage)
	server.Shutdown(ctx)

}
