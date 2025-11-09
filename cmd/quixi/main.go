package main

import (
	"context"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/config"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/devserver"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.NewConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	quixi := NewService(conf)
	if quixi.conf.EnableDevServer {
		// Запуск фоновых сервисов
		go quixi.DevService.Start(ctx)
		// Запуск HTTP-сервера (он должен блокировать!)
		httpServer := devserver.NewHTTPServer(conf.DevServerWeb)
		go func() {
			httpServer.Start()
		}()
		log.Println("Dev mode active")
	} else {
		if err := quixi.Listener.Connect(); err != nil {
			log.Fatal("Failed to connect:", err)
		}
		defer quixi.Listener.Disconnect()
		quixi.Listener.GetMessages(ctx)
		quixi.Listener.Subscribe("#", handleMessage)
		log.Println("MQTT Broker running...")
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // Wait for interrupt signal
	log.Println("Shutting down gracefully...")
	gracefulShutdown(cancel)
}