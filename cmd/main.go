package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	service "github.com/onnenon/whiteboard_go/internal/service"

	nats "github.com/nats-io/nats.go"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	ws = service.WhiteboardService{}

	go waitForInterrupt(&wg, nc)

	wg.Wait()
}

func waitForInterrupt(wg *sync.WaitGroup, nc *nats.Conn) {
	defer wg.Done()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	sig := <-ch
	log.Printf("Caught signal: %s - draining\n", sig)
	nc.Drain()
}
