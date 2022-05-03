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

	registerSignalShutdown(stop)
	defer registerRecover(stop)
	app, err := internal.NewApplication(*configDir, ctx)
	if err != nil {
		log.Fatalln(err)
	}
	app.Run()
}

func registerSignalShutdown(stop context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		log.Println("Got system stop signal")
		stop()
	}()
}

func registerRecover(stop context.CancelFunc) {
	if err := recover(); err != nil {
		log.Println("Panic! Error: \n", err)
		stop()
	}
}
