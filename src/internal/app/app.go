package app

import (
	"log"

	"github.com/lekht/account-master/src/config"
)

func Run(cfg *config.Config) {
	log.Printf("config: %+v\n", *cfg)
}
