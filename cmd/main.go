package main

import (
	"house-service/config"
	"house-service/internal/app"
	"log"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal("can't read config file")
	}

	app.Run(cfg)
}