package main

import (
	"log"

	"github.com/logica0419/scheduled-messenger-bot/config"
	"github.com/logica0419/scheduled-messenger-bot/repository"
	"github.com/logica0419/scheduled-messenger-bot/router"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

func main() {
	log.Print("Initializing Scheduled Messenger Bot...")

	c, err := config.GetConfig()
	if err != nil {
		log.Panicf("Error: failed to get config - %s", err)
	}

	api := api.GetApi(c)

	repo, err := repository.GetRepository(c)
	if err != nil {
		log.Panicf("Error: failed to initialize DB - %s", err)
	}

	r := router.SetUpRouter(c, api, repo)

	r.Start()
}
