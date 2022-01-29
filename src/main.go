package main

import (
	"log"
	"versioncontrol-service/src/config"
	"versioncontrol-service/src/server"
	"versioncontrol-service/src/services"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal("Unable to load configurations", err)
	}
	services.InitializeOAuthGithub()
	server.Start()
}
