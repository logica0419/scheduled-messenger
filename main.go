package main

import (
	"fmt"

	"github.com/logica0419/scheduled-messenger-bot/router"
)

func main() {
	fmt.Println("Initializing Scheduled Messenger Bot...")

	e := router.Setup()

	e.Logger.Panic(e.Start(":8080"))
}
