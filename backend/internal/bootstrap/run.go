// Package bootstrap wires dependencies and coordinates application shutdown.
package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"logmanagement/backend/internal/collector"
	"logmanagement/backend/internal/config"
	"logmanagement/backend/internal/repository"
	"logmanagement/backend/internal/retention"
	"logmanagement/backend/internal/router"
)

func Run(parent context.Context, cfg config.Config) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	initCtx, cancelInit := context.WithTimeout(ctx, 30*time.Second)
	defer cancelInit()
	store, err := repository.OpenPostgres(initCtx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database startup failed: %w", err)
	}
	defer store.Pool.Close()
	if err = store.Seed(initCtx, cfg.AdminPassword, cfg.ViewerPassword, cfg.APIKeyA, cfg.APIKeyB); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}
	receiver, err := collector.Open(initCtx, store, cfg.SyslogAddress, cfg.SyslogTenant)
	if err != nil {
		return err
	}
	defer receiver.Close()
	listener, err := net.Listen("tcp", cfg.HTTPAddress)
	if err != nil {
		return fmt.Errorf("listen HTTP: %w", err)
	}
	defer listener.Close()
	cancelInit()

	server := router.NewServer(store, router.Config{Origin: cfg.Origin, SecureCookies: cfg.SecureCookies})
	var workers sync.WaitGroup
	workers.Add(2)
	go func() { defer workers.Done(); receiver.Run(ctx) }()
	go func() { defer workers.Done(); retention.Run(ctx, store, cfg.RetentionDays) }()
	// Workers must exit before the database pool is closed, including on HTTP failure.
	defer func() { cancel(); workers.Wait() }()

	shutdownDone := make(chan struct{})
	context.AfterFunc(ctx, func() {
		defer close(shutdownDone)
		if err := server.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Printf("HTTP shutdown: %v", err)
		}
		// Also handles cancellation before Fiber has attached this pre-bound listener.
		listener.Close()
	})
	log.Printf("API %s; Syslog UDP/TCP %s (%s)", cfg.HTTPAddress, cfg.SyslogAddress, cfg.SyslogTenant)
	err = server.Listener(listener)
	canceled := ctx.Err() != nil
	cancel()
	<-shutdownDone
	if err != nil && !canceled {
		return fmt.Errorf("HTTP server: %w", err)
	}
	return nil
}
