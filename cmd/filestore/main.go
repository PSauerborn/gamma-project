package main

import "github.com/PSauerborn/gamma-project/pkg/utils"

var cfg = utils.NewConfigMapWithValues(map[string]string{
	"listen_address": "0.0.0.0",
	"listen_port":    "10311",
	"postgres_url":   "postgres://postgres:postgres-dev@192.168.99.100:5432/gamma_project",
})

func main() {
	cfg.ConfigureLogging()

}
