package main

import (
	"TwitchTokGen/config"
	"TwitchTokGen/webserver"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nicklaw5/helix/v2"
)

func main() {
	config.Load("./settings.json")
	log.Println("Start Token Gateway")
	go webserver.Start(&helix.Options{
		ClientID:     config.Config.Twitch.Api.ClientId,
		ClientSecret: config.Config.Twitch.Api.ClientSecret,
		RedirectURI:  config.Config.Twitch.Api.RedirectURL,
	})

	// This prevents the process from exiting
	fmt.Println("Software is running. Press CTRL-C to exit.")
	exchan := make(chan os.Signal, 1)
	signal.Notify(exchan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-exchan
	log.Println("Exit")
}
