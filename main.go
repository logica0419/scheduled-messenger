package main

import (
	"log"

	"github.com/logica0419/scheduled-messenger-bot/config"
	"github.com/logica0419/scheduled-messenger-bot/router"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

func main() {
	log.Print("Initializing Scheduled Messenger Bot...")

	if err := config.GetConfig(); err != nil {
		log.Panic(err)
	}

	api.SetUpApi()

	r := router.Setup()

	r.Start()
}
