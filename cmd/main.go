package main

import (
	"context"
	application "github.com/bifidokk/recipe-bot/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	app, err := application.NewApp(ctx)

	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = app.Run()

	if err != nil {
		log.Fatalf("failed to run: %s", err.Error())
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	log.Println("Shutdown signal received, exiting...")
}
