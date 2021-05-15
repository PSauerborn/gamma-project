package main

import (
	"fmt"
	"strconv"

	"github.com/PSauerborn/gamma-project/pkg/filestore"
	"github.com/PSauerborn/gamma-project/pkg/utils"
)

var cfg = utils.NewConfigMapWithValues(map[string]string{
	"listen_port":  "10314",
	"postgres_url": "postgres://postgres:development@localhost:5432/gamma_project",
	"log_level":    "DEBUG",
	"storage_path": "./data",
})

func main() {
	cfg.ConfigureLogging()
	// get listen port from env vars and convert to int
	port, err := strconv.Atoi(cfg.Get("listen_port"))
	if err != nil {
		panic(fmt.Sprintf("received invalid port '%s'", cfg.Get("listen_port")))
	}

	// generate new persistence layer and connect
	persistence := filestore.NewPostgresPersistence(cfg.Get("postgres_url"),
		cfg.Get("storage_path"))
	if err := persistence.Connect(); err != nil {
		panic(fmt.Errorf("unable to connect persistence: %+v", err))
	}
	defer persistence.Close()
	// generate new instance of API and run
	filestore.NewFilestoreAPI(persistence).Run(fmt.Sprintf(":%d", port))
}
