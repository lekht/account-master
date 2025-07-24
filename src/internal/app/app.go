package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lekht/account-master/src/config"
	"github.com/lekht/account-master/src/internal/controllers"
	"github.com/lekht/account-master/src/internal/hash"
	"github.com/lekht/account-master/src/internal/model"
	"github.com/lekht/account-master/src/pkg/server"
	"github.com/lekht/account-master/src/pkg/storage/mock"
)

func Run(cfg *config.Config) {
	log.Printf("config: %+v\n", *cfg)

	storage := mock.New()

	// create admin
	{
		hash, err := hash.HashPassword(cfg.Admin.Password)
		if err != nil {
			log.Panicf("failed to hash admin pwd: %v\n", err)
		}

		err = storage.CreateUser(model.Profile{
			Email:    cfg.Admin.Email,
			Username: cfg.Admin.Username,
			Password: hash,
			Admin:    cfg.Admin.Admin,
		})
		if err != nil {
			log.Panicf("failed to create admin: %v\n", err)
		}
	}

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
