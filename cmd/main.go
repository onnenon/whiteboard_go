package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	nats "github.com/nats-io/nats.go"
)

func statusHander(m *nats.Msg) {
	fmt.Printf("Message %s\n", string(m.Data))
}

func disconnectHandler(nc *nats.Conn) {
	log.Printf("Disconnected from %s\n", nc.Servers())
}

func reconnectHandler(nc *nats.Conn) {
	log.Printf("Reconnected to %s\n", nc.ConnectedUrl())
}

func getOptions() []nats.Option {
	opt := []nats.Option{}
	opt = append(opt, nats.DisconnectHandler(disconnectHandler))
	opt = append(opt, nats.ReconnectHandler(reconnectHandler))
	return opt
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	opts := getOptions()

	nc, err := nats.Connect("0.0.0.0:4444", opts...)

	if err != nil {
		log.Fatal((err))
	}

	fmt.Printf("Connected to %s\n", nc.ConnectedUrl())

	nc.Subscribe("status", statusHander)

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
