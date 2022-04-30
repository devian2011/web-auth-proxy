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

var (
	configFile = flag.String("config", "./config/main.config.json", "Configuration file path")
)

func main() {
	flag.Parse()
	_ = flag.Lookup("log_dir").Value.Set("./logs")
	_ = flag.Lookup("alsologtostderr").Value.Set("true")

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
			stop()
		}
	}()

	log.Println("Application init...")
	app := internal.NewApplication(*configFile, ctx)
	log.Println("Application is starting...")
	app.Run()
	log.Println("Application is running...")
}

