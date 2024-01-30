package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/MrColorado/backend/bookHandler/internal/config"
	"github.com/nats-io/nats.go"
)

type msgHandler func([]byte, string)
type requestHandler func([]byte, string) ([]byte, error)

type subDataHolder struct {
	input chan *nats.Msg
	sub   *nats.Subscription
}

type NatsClient struct {
	wg  sync.WaitGroup
	ctx context.Context

	conn       *nats.Conn
	input      chan *nats.Msg
	subDataMap map[string]subDataHolder
}

func NewNatsClient(cfg config.NatsConfigStruct, ctx context.Context) (*NatsClient, error) {
	fmt.Printf("Connect nats at : %s\n", cfg.NatsHOST)
	conn, err := nats.Connect(cfg.NatsHOST, nats.Name("bookHandler"))
	if err != nil {
		fmt.Printf("Failed to connect nats at : %s\n", cfg.NatsHOST)
		fmt.Println(err.Error())
		return nil, fmt.Errorf("failed to create nats client")
	}
	fmt.Printf("Succeed to connect nats at : %s\n", cfg.NatsHOST)
	return &NatsClient{
		wg:         sync.WaitGroup{},
		ctx:        ctx,
		conn:       conn,
		subDataMap: make(map[string]subDataHolder),
	}, nil
}

func (nc *NatsClient) AddChanQueueSub(subject string, group string) error {
	fmt.Printf("Listen on %s with group : %s\n", subject, group)
	input := make(chan *nats.Msg)
	sub, err := nc.conn.ChanQueueSubscribe(subject, group, input)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to subscribe to %s with group %s", subject, group)
	}
	fmt.Printf("Succeed to listen on %s with group : %s\n", subject, group)
	nc.subDataMap[subject] = subDataHolder{
		input: input,
		sub:   sub,
	}
	return nil
}

func (nc *NatsClient) RemoveChanQueueSub(subject string) error {
	holder, ok := nc.subDataMap[subject]
	if !ok {
		return fmt.Errorf("service is not listening on %s", subject)
	}

	err := holder.sub.Unsubscribe()
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to unsubcribe from %s", subject)
	}

	close(holder.input)
	delete(nc.subDataMap, subject)
	return nil
}

func (nc *NatsClient) Run(msgHdl msgHandler, rqthandler requestHandler) {
	fmt.Printf("Run")
	for _, holder := range nc.subDataMap {
		nc.wg.Add(1)
		go func(nc *NatsClient, holder subDataHolder) {
			fmt.Printf("Listen on : %s on queue : %s", holder.sub.Subject, holder.sub.Queue)
			defer nc.wg.Done()
			for {
				select {
				case msg, ok := <-holder.input:
					if !ok {
						fmt.Printf("KO Listen on : %s on queue : %s", holder.sub.Subject, holder.sub.Queue)
						return
					}
					fmt.Printf("Got msg : %s ", msg.Subject)
					nc.input <- msg
				case <-nc.ctx.Done():
					fmt.Printf("DONE Listen on : %s on queue : %s", holder.sub.Subject, holder.sub.Queue)
					return
				}
			}
		}(nc, holder)
		fmt.Printf("Run end")
	}

out:
	for {
		select {
		case msg := <-nc.input:
			if msg == nil {
				fmt.Println("msg received is nil")
				continue
			}
			if msg.Reply == "" {
				msgHdl(msg.Data, msg.Subject)
			} else {
				rsp, err := rqthandler(msg.Data, msg.Subject)
				if err != nil {
					fmt.Println(err.Error())
				}
				nc.conn.Publish(msg.Reply, rsp)
			}
		case <-nc.ctx.Done():
			break out
		}
	}

	nc.wg.Wait()
}
