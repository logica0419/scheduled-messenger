package main

import (
	"log"

	"github.com/logica0419/scheduled-messenger-bot/config"
	"github.com/logica0419/scheduled-messenger-bot/router"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

func main() {
	log.Print("Initializing Scheduled Messenger Bot...")

	c, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	api := api.GetApi(c)

	r := router.SetUpRouter(c, api)

	r.Start()
}
