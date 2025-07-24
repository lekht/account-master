package main

import (
	"flag"
	"log"

	"github.com/lekht/account-master/src/config"
	"github.com/lekht/account-master/src/docs"
	"github.com/lekht/account-master/src/internal/app"
)

var conf config.Config

// @title						Account Master
// @version					1.0
// @decsription				CRUD account service
// @BasePath					/
// @securityDefinitions.basic	BasicAuth
func main() {
	docs.SwaggerInfo.BasePath = "/"

	path := flag.String("config", "", "path to config file")
	flag.Parse()

	if *path != "" {
		if err := config.Load(*path, &conf); err != nil {
			log.Fatalf("failed to load config: %v\n", err)
		}
	} else {
		if err := config.Load("./config.yaml", &conf); err != nil {
			log.Fatalf("failed to load config: %v\n", err)
		}
	}

	app.Run(&conf)
}
