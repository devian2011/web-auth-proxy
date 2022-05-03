package main

import (
	"context"
	"flag"
	"log"
	"lproxy/internal"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configDir := flag.String("configDir", "./config", "Configuration directory")
	flag.Parse()
	ctx, stop := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		log.Println("Got system stop signal")
		stop()
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic! Error: \n", err)
			stop()
		}
	}()

	log.Println("Application init...")
	app := internal.NewApplication(*configDir, ctx)
	log.Println("Application is starting...")
	app.Run()
	log.Println("Application is running...")
}

