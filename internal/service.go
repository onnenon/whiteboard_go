package service

import (
	"fmt"
	"log"

	nats "github.com/nats-io/nats.go"
)

// Whiteboard Service
type WhiteboardService struct {
	conn *nats.Conn
}

func statusHander(m *nats.Msg) {
	fmt.Printf("Message %s\n", string(m.Data))
}

func disconnectHandler(nc *nats.Conn) {
	log.Printf("Disconnected from %s\n", nc.Servers())
}

func reconnectHandler(nc *nats.Conn) {
	log.Printf("Reconnected to %s\n", nc.ConnectedUrl())
}

func (ws *WhiteboardService) getOptions() []nats.Option {
	opt := []nats.Option{}
	opt = append(opt, nats.DisconnectHandler(disconnectHandler))
	opt = append(opt, nats.ReconnectHandler(reconnectHandler))
	return opt
}

func (ws *WhiteboardService) init() {
	var err error
	ws.conn, err = nats.Connect("0.0.0.0:4444", ws.getOptions()...)

	if err != nil {
		log.Fatal((err))
	}

	fmt.Printf("Connected to %s\n", ws.conn.ConnectedUrl())

	ws.conn.Subscribe("status", statusHander)
}
