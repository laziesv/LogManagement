package collector

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"logmanagement/backend/internal/model"
)

type testStore struct {
	events  chan model.Event
	ruleErr error
}

func (s *testStore) GetRule(context.Context, string) (model.Rule, error) {
	return model.Rule{}, s.ruleErr
}
func (s *testStore) Insert(ctx context.Context, events []model.Event) error {
	for _, e := range events {
		select {
		case s.events <- e:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func TestUDPAndTCPIngestAndShutdown(t *testing.T) {
	store := &testStore{events: make(chan model.Event, 4)}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, err := Open(ctx, store, "127.0.0.1:0", "demo-a")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	done := make(chan struct{})
	go func() { defer close(done); r.Run(ctx) }()
	udp, err := net.Dial("udp", r.udp.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer udp.Close()
	tcp, err := net.DialTimeout("tcp", r.tcp.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer tcp.Close()
	line := "<134>fw01 product=ngfw action=deny src=10.0.1.10 tenant=demo-b\n"
	if _, err = udp.Write([]byte(line)); err != nil {
		t.Fatal(err)
	}
	if _, err = tcp.Write([]byte(line)); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		select {
		case e := <-store.events:
			if e.Tenant != "demo-a" || e.Source != "firewall" || e.SrcIP != "10.0.1.10" {
				t.Fatalf("bad normalized event: %+v", e)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("missing Syslog event")
		}
	}
	// Leave TCP idle: shutdown must close active clients, not wait for the 30s deadline.
	cancel()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("collector failed to stop with an idle TCP connection")
	}
	tcp.SetReadDeadline(time.Now().Add(time.Second))
	if _, err := tcp.Read(make([]byte, 1)); err == nil {
		t.Fatal("TCP connection left open")
	} else if e, ok := err.(net.Error); ok && e.Timeout() {
		t.Fatal("TCP read timed out instead of seeing closed connection")
	}
}

func TestRejectInvalidTenant(t *testing.T) {
	store := &testStore{ruleErr: errors.New("missing tenant")}
	if r, err := Open(context.Background(), store, "127.0.0.1:0", "missing"); err == nil {
		r.Close()
		t.Fatal("accepted invalid tenant")
	}
}

func TestTCPBindFailureReleasesUDP(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer occupied.Close()
	address := occupied.Addr().String()
	if r, err := Open(context.Background(), &testStore{}, address, "demo-a"); err == nil {
		r.Close()
		t.Fatal("expected TCP bind error")
	}
	udp, err := net.ListenPacket("udp", address)
	if err != nil {
		t.Fatalf("UDP socket leaked on startup failure: %v", err)
	}
	udp.Close()
}
