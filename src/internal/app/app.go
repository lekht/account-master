package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lekht/account-master/src/config"
	"github.com/lekht/account-master/src/internal/controllers"
	"github.com/lekht/account-master/src/pkg/server"
	"github.com/lekht/account-master/src/pkg/storage"
)

func Run(cfg *config.Config) {
	log.Printf("config: %+v\n", *cfg)

	storage, err := storage.New()
	if err != nil {
		log.Fatalf("failed to create new storage: %v", err)
	}

	router := controllers.New(storage)

	httpserver := server.New(router.GetRouter(), server.Adress(cfg.Server.Host, cfg.Server.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err := <-httpserver.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpserver.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
