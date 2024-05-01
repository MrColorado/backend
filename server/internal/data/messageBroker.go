package data

import (
	"time"

	"github.com/MrColorado/backend/logger"
	"github.com/MrColorado/backend/server/internal/config"
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	conn *nats.Conn
}

func NewNatsClient(cfg config.NatsConfigStruct) (*NatsClient, error) {
	logger.Infof("Connect to nats at : %s", cfg.NatsHOST)
	conn, err := nats.Connect(cfg.NatsHOST)
	if err != nil {
		return nil, logger.Errorf("failed to create nats client")
	}
	return &NatsClient{
		conn: conn,
	}, nil
}

func (nc *NatsClient) PublishMsg(topic string, msg []byte) error {
	logger.Infof("Publish on : %s", topic)
	err := nc.conn.Publish(topic, msg)
	if err != nil {
		return logger.Errorf("failed to publish %s on %s", msg, topic)
	}
	return nil
}

func (nc *NatsClient) Request(topic string, msg []byte) ([]byte, error) {
	logger.Infof("Request on : %s", topic)
	rsp, err := nc.conn.Request(topic, msg, 5*time.Second)
	if err != nil {
		return []byte{}, logger.Errorf("failed to request %s on %s", msg, topic)
	}
	return rsp.Data, nil
}
