package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lekht/account-master/src/config"
	"github.com/lekht/account-master/src/internal/controllers"
	"github.com/lekht/account-master/src/internal/model"
	"github.com/lekht/account-master/src/pkg/server"
	"github.com/lekht/account-master/src/pkg/storage/mock"
)

func Run(cfg *config.Config) {
	log.Printf("config: %+v\n", *cfg)

	storage := mock.New()

	// create main superuser with id=0
	storage.CreateUser(model.Profile{
		Email:    "a@mail.com",
		Username: "admin",
		Password: "aaa",
		Admin:    true,
	})

	router := controllers.New(storage)

	httpserver := server.New(router.Router(), server.Adress(cfg.Server.Host, cfg.Server.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
		break
	case err := <-httpserver.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err := httpserver.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
