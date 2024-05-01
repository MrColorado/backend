package handler

import (
	"context"
	"sync"

	"github.com/MrColorado/backend/book-handler/internal/config"
	"github.com/MrColorado/backend/logger"
	"github.com/nats-io/nats.go"
)

type msgHandler func([]byte, string)
type requestHandler func([]byte, string) ([]byte, error)

type subDataHolder struct {
	input chan *nats.Msg
	sub   *nats.Subscription
}

type NatsClient struct {
	ctx        context.Context
	conn       *nats.Conn
	input      chan *nats.Msg
	subDataMap map[string]subDataHolder
	wg         sync.WaitGroup
}

func NewNatsClient(cfg config.NatsConfigStruct, ctx context.Context) (*NatsClient, error) {
	logger.Infof("Connect nats at : %s", cfg.NatsHOST)
	conn, err := nats.Connect(cfg.NatsHOST, nats.Name("book-handler"))
	if err != nil {
		logger.Info(err.Error())
		return nil, logger.Error("failed to create nats client")
	}
	logger.Infof("Succeed to connect nats at : %s", cfg.NatsHOST)
	return &NatsClient{
		wg:         sync.WaitGroup{},
		ctx:        ctx,
		conn:       conn,
		input:      make(chan *nats.Msg),
		subDataMap: make(map[string]subDataHolder),
	}, nil
}

func (nc *NatsClient) AddChanQueueSub(subject string, group string) error {
	logger.Infof("Listen on %s with group : %s", subject, group)

	input := make(chan *nats.Msg)
	sub, err := nc.conn.ChanQueueSubscribe(subject, group, input)
	if err != nil {
		logger.Info(err.Error())
		return logger.Errorf("failed to subscribe to %s with group %s", subject, group)
	}
	logger.Infof("Succeed to listen on %s with group : %s", subject, group)
	nc.subDataMap[subject] = subDataHolder{
		input: input,
		sub:   sub,
	}

	nc.wg.Add(1)
	go func(nc *NatsClient, holder subDataHolder) {
		logger.Infof("Listen on : %s on queue : %s", holder.sub.Subject, holder.sub.Queue)
		defer nc.wg.Done()
		for {
			select {
			case msg, ok := <-holder.input:
				if !ok {
					logger.Infof("Ko Listen on : %s on queue : %s", holder.sub.Subject, holder.sub.Queue)
					return
				}
				logger.Infof("Got msg %s on %s", string(msg.Data), msg.Subject)
				nc.input <- msg
			case <-nc.ctx.Done():
				logger.Infof("Done listen on : %s on queue : %s", holder.sub.Subject, holder.sub.Queue)
				return
			}
		}
	}(nc, nc.subDataMap[subject])

	return nil
}

func (nc *NatsClient) RemoveChanQueueSub(subject string, del bool) error {
	logger.Infof("Remove sub on %s", subject)

	holder, ok := nc.subDataMap[subject]
	if !ok {
		return logger.Errorf("service is not listening on %s", subject)
	}

	err := holder.sub.Unsubscribe()
	if err != nil {
		logger.Info(err.Error())
		return logger.Errorf("failed to unsubcribe from %s", subject)
	}

	close(holder.input)
	if del {
		delete(nc.subDataMap, subject)
	}

	return nil
}

func (nc *NatsClient) Run(msgHdl msgHandler, rqthandler requestHandler) {
	logger.Info("Run")

out:
	for {
		select {
		case msg := <-nc.input:
			if msg == nil {
				logger.Info("msg received is nil")
				continue
			}
			logger.Infof("received on %s msg : ", msg.Subject, string(msg.Data))
			if msg.Reply == "" {
				msgHdl(msg.Data, msg.Subject)
			} else {
				rsp, err := rqthandler(msg.Data, msg.Subject)
				if err != nil {
					logger.Info(err.Error())
				}
				nc.conn.Publish(msg.Reply, rsp)
			}
		case <-nc.ctx.Done():
			break out
		}
	}

	nc.wg.Wait()
}

// TODO : AddChanQueueSub wiil not read if called after run
