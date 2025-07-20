package main

import (
	"flag"
	"log"

	"github.com/lekht/account-master/src/config"
	"github.com/lekht/account-master/src/internal/app"
)

var conf config.Config

func init() {
	path := flag.String("config", "", "path to config file")
	flag.Parse()

	err := config.Load(*path, &conf)
	if err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}
}

func main() {
	app.Run(&conf)
}
