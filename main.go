package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var rateLimit struct {
	RateLimit_Limit     int
	RateLimit_Remaining int
	RateLimit_Reset     int
}

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

	log.Println("Start Token Gateway", config)
	go startWebserver()

	go func() {
		log.Println("Starting Ratelimit Handler", config.HttpRequest.RateLimit)
		if config.HttpRequest.RateLimit > 0 {

			rateLimit.RateLimit_Limit = config.HttpRequest.RateLimit
			rateLimit.RateLimit_Remaining = config.HttpRequest.RateLimit
			rateLimit.RateLimit_Reset = 0
			for {
				if rateLimit.RateLimit_Remaining < rateLimit.RateLimit_Limit {
					rateLimit.RateLimit_Remaining += 1
					rateLimit.RateLimit_Reset = ((rateLimit.RateLimit_Limit - rateLimit.RateLimit_Remaining) / config.HttpRequest.RegenerateRate)
					log.Printf("\nLimit: %d\nRemaining: %d\nRegenerate: %d\n", rateLimit.RateLimit_Limit, rateLimit.RateLimit_Remaining, rateLimit.RateLimit_Reset)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

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
