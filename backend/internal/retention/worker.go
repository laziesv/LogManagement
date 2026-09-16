// Package retention owns scheduled log/session cleanup.
package retention

import (
	"context"
	"log"
	"time"
)

type Cleaner interface {
	Cleanup(context.Context, int) error
}

const interval = 1 * time.Minute

// Run blocks until cancellation. Cleanup runs immediately and then every 1 minute.
func Run(ctx context.Context, store Cleaner, days int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if ctx.Err() != nil {
			return
		}
		cleanupCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		if err := store.Cleanup(cleanupCtx, days); err != nil && ctx.Err() == nil {
			log.Printf("cleanup: %v", err)
		}
		cancel()
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
