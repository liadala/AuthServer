package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		err         error
		exitChannel chan os.Signal
		//config => settings.go
	)

	config, err = loadConfig("./settings.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start Token Gateway")
	go startWebserver()

	// Wait for Terminate Process
	log.Println("press CTRL-C to exit")
	exitChannel = make(chan os.Signal, 1)
	signal.Notify(exitChannel, syscall.SIGTERM, syscall.SIGINT)
	signal := <-exitChannel
	if signal == syscall.SIGTERM {
		log.Println("received SIGTERM signal")
	} else if signal == syscall.SIGINT {
		log.Println("received SIGINT signal")
	}
}
