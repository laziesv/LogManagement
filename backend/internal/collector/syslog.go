// Package collector owns the UDP/TCP Syslog receiver lifecycle.
package collector

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"logmanagement/backend/internal/model"
	"logmanagement/backend/internal/normalize"
)

type Store interface {
	GetRule(context.Context, string) (model.Rule, error)
	Insert(context.Context, []model.Event) error
}

type Receiver struct {
	store  Store
	tenant string
	udp    net.PacketConn
	tcp    net.Listener
}

// Open validates the deployment tenant and binds both protocols before starting work.
func Open(ctx context.Context, store Store, address, tenant string) (*Receiver, error) {
	if _, err := store.GetRule(ctx, tenant); err != nil {
		return nil, fmt.Errorf("invalid SYSLOG_TENANT: %w", err)
	}
	udp, err := net.ListenPacket("udp", address)
	if err != nil {
		return nil, fmt.Errorf("listen Syslog UDP: %w", err)
	}
	tcp, err := net.Listen("tcp", address)
	if err != nil {
		udp.Close()
		return nil, fmt.Errorf("listen Syslog TCP: %w", err)
	}
	return &Receiver{store: store, tenant: tenant, udp: udp, tcp: tcp}, nil
}

// Close releases listening sockets. Cancel Run's context to stop active clients too.
func (r *Receiver) Close() { r.udp.Close(); r.tcp.Close() }

// Run waits for listeners and active clients to exit on context cancellation.
func (r *Receiver) Run(ctx context.Context) {
	stopClose := context.AfterFunc(ctx, r.Close)
	defer stopClose()
	defer r.Close()
	var workers sync.WaitGroup
	workers.Add(2)
	go func() { defer workers.Done(); r.receiveUDP(ctx) }()
	go func() { defer workers.Done(); r.receiveTCP(ctx) }()
	workers.Wait()
}

func (r *Receiver) ingest(ctx context.Context, line string) {
	job, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	event, err := normalize.NormalizeSyslog(line, r.tenant, time.Now().UTC())
	if err == nil {
		err = r.store.Insert(job, []model.Event{event})
	}
	if err != nil && ctx.Err() == nil {
		log.Printf("syslog rejected: %v", err)
	}
}

func (r *Receiver) receiveUDP(ctx context.Context) {
	buf := make([]byte, 65536)
	for {
		n, _, err := r.udp.ReadFrom(buf)
		if err != nil {
			return
		}
		r.ingest(ctx, string(buf[:n]))
	}
}

func (r *Receiver) receiveTCP(ctx context.Context) {
	slots := make(chan struct{}, 32)
	var clients sync.WaitGroup
	defer clients.Wait()
	for {
		conn, err := r.tcp.Accept()
		if err != nil {
			return
		}
		select {
		case slots <- struct{}{}:
		default:
			conn.Close()
			continue
		}
		clients.Add(1)
		go func() {
			defer clients.Done()
			defer func() { <-slots }()
			r.handleTCP(ctx, conn)
		}()
	}
}

func (r *Receiver) handleTCP(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	stopClose := context.AfterFunc(ctx, func() { conn.Close() })
	defer stopClose()
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 4096), 65536)
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	for scanner.Scan() {
		r.ingest(ctx, scanner.Text())
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	}
}
