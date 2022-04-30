package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"lproxy/internal"
)

var (
	configFile = flag.String("config", "./config/main.config.json", "Configuration file path")
)

func main() {
	flag.Parse()
	_ = flag.Lookup("log_dir").Value.Set("./logs")
	_ = flag.Lookup("alsologtostderr").Value.Set("true")

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		log.Println("Got system stop signal")
		done <- true
	}()
	log.Println("Application init...")
	app := internal.NewApplication(*configFile)
	log.Println("Application is starting...")
	app.Run()
	log.Println("Application is running...")

	<-done
	app.Stop()
	log.Println("Application has been stopped successfully")
}

