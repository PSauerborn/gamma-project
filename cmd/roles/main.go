package main

import (
	"fmt"
	"strconv"

	"github.com/PSauerborn/gamma-project/pkg/roles"
	"github.com/PSauerborn/gamma-project/pkg/utils"
)

var cfg = utils.NewConfigMapWithValues(map[string]string{
	"listen_address": "0.0.0.0",
	"listen_port":    "10313",
	"postgres_url":   "postgres://postgres:postgres-dev@192.168.99.100:5432/gamma_project",
})

func main() {
	cfg.ConfigureLogging()
	// generate new postgres persistence and connect
	db := roles.NewPostgresPersistence(cfg.Get("postgres_url"))
	if err := db.Connect(); err != nil {
		panic(fmt.Errorf("unable to connect postgres persistence: %+v", err))
	}
	defer db.Close()

	// parse listen port into integer
	listenPort, err := strconv.Atoi(cfg.Get("listen_port"))
	if err != nil {
		panic(fmt.Sprintf("received invalid listen port %s", cfg.Get("listen_port")))
	}

	roles.NewRolesAPI(db).Run(fmt.Sprintf(":%d", listenPort))
}
