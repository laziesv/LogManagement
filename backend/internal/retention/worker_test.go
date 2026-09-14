package retention

import (
	"context"
	"testing"
	"time"
)

type cleanupFunc func(context.Context, int) error

func (f cleanupFunc) Cleanup(ctx context.Context, days int) error { return f(ctx, days) }

func TestImmediateCleanupAndCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	called := make(chan struct{})
	done := make(chan struct{})
	cleaner := cleanupFunc(func(job context.Context, days int) error {
		if days != 14 {
			t.Errorf("days = %d", days)
		}
		deadline, ok := job.Deadline()
		if !ok || time.Until(deadline) > 30*time.Second {
			t.Error("cleanup lacks bounded deadline")
		}
		close(called)
		<-job.Done()
		return job.Err()
	})
	go func() { defer close(done); Run(ctx, cleaner, 14) }()
	select {
	case <-called:
	case <-time.After(3 * time.Second):
		t.Fatal("initial cleanup did not run")
	}
	cancel()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("cleanup did not stop on cancellation")
	}
}

func TestAlreadyCanceledSkipsCleanup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	Run(ctx, cleanupFunc(func(context.Context, int) error { t.Error("cleanup ran after cancellation"); return nil }), 7)
}
