package dataHandler

import (
	"fmt"
	"time"

	"github.com/MrColorado/backend/server/internal/config"
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	conn *nats.Conn
}

func NewNatsClient(cfg config.NatsConfigStruct) (*NatsClient, error) {
	fmt.Printf("Connect to nats at : %s\n", cfg.NatsHOST)
	conn, err := nats.Connect(cfg.NatsHOST)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("failed to create nats client")
	}
	return &NatsClient{
		conn: conn,
	}, nil
}

func (nc *NatsClient) PublishMsg(topic string, msg []byte) error {
	err := nc.conn.Publish(topic, msg)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to publish %s on %s", msg, topic)
	}
	return nil
}

func (nc *NatsClient) Request(topic string, msg []byte) ([]byte, error) {
	rsp, err := nc.conn.Request(topic, msg, 10*time.Microsecond)
	if err != nil {
		fmt.Println(err.Error())
		return []byte{}, fmt.Errorf("failed to request %s on %s", msg, topic)
	}
	return rsp.Data, nil
}
