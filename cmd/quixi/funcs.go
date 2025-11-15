package main

import (
	"context"
	"log"
)
func gracefulShutdown(cancel context.CancelFunc) {
	cancel()
	log.Println("All workers stopped gracefully")
}
