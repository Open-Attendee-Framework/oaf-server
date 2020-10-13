package main

import (
	"log"

	"github.com/concertLabs/oaf-server/internal/packer"
	"github.com/concertLabs/oaf-server/pkg/config"
	db100 "github.com/concertLabs/oaf-server/pkg/db/v100"
)

func main() {
	packer.PackAll()
	conf, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	db100.Initialisation(conf.DatabaseConnection)
}
