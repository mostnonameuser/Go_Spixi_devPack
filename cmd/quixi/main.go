package main

import (
	"context"
	"fmt"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/config"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/devserver"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/network"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.NewConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var quixi network.Broker
	if conf.EnableDevServer {
		httpServer := devserver.NewHTTPServer(conf.DevServerWeb)
		go func() {
			httpServer.Start()
		}()
		quixi = network.NewWsBroker(conf)
	} else {
		quixi = network.NewMQQTBroker(conf)
	}
	if err := quixi.Connect(); err != nil {
		log.Fatal("Failed to connect broker:", err)
	}
	defer quixi.Disconnect()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		fmt.Println("Starting server")
		if err := quixi.Start(ctx); err != nil {
			log.Fatal("Broker failed:", err)
		}
	}()
	<-sigChan
	log.Println("Shutting down gracefully...")

	if err := quixi.Disconnect(); err != nil {
		log.Printf("Error during disconnect: %v", err)
	}
	gracefulShutdown(cancel)
}
