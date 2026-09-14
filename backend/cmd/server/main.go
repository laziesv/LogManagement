package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"logmanagement/backend/internal/bootstrap"
	"logmanagement/backend/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return bootstrap.Run(ctx, cfg)
}
