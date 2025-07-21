package app

import (
	"log"

	"github.com/lekht/account-master/src/config"
	"github.com/lekht/account-master/src/pkg/storage"
)

func Run(cfg *config.Config) {
	log.Printf("config: %+v\n", *cfg)

	_, err := storage.New()
	if err != nil {
		log.Fatalf("failed to create new storage: %v", err)
	}
}
