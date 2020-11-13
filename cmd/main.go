package main

import (
	"fmt"
	"log"

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
	opts := getOptions()

	nc, err := nats.Connect("0.0.0.0:4444", opts...)

	if err != nil {
		log.Fatal((err))
	}

	fmt.Printf("Connected to %s\n", nc.ConnectedUrl())

	defer nc.Close()

	nc.Subscribe("status", statusHander)
	for {
	}
}
