package nats_listener

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/nats-io/stan.go"
	"github.com/vanamelnik/wildberries-L0/models"
)

type NATSListener struct {
	sc  stan.Conn
	sub stan.Subscription
}

func New(stanCluster, clientID, durableName, subject string) (NATSListener, error) {
	sc, err := stan.Connect(stanCluster, clientID)
	if err != nil {
		return NATSListener{}, err
	}
	nl := NATSListener{
		sc: sc,
	}
	sub, err := sc.Subscribe(subject, nl.msgHandler, stan.DurableName(durableName))
	if err != nil {
		return NATSListener{}, err
	}
	nl.sub = sub

	return nl, nil
}

func (nl NATSListener) Close() (retErr error) {
	if err := nl.sub.Close(); err != nil {
		retErr = multierror.Append(retErr, err)
	}
	if err := nl.sc.Close(); err != nil {
		retErr = multierror.Append(retErr, err)
	}
	return
}

func (nl NATSListener) msgHandler(msg *stan.Msg) {
	var order models.Order
	if err := json.Unmarshal(msg.Data, &order); err != nil {
		log.Printf("natsListener: ERR: order rejected: incorrect order type: %s", err)
		return
	}
	if err := order.Validate(); err != nil {
		log.Printf("natsListener: ERR: order rejected: invalid order: %s", err)
		return
	}
	log.Printf("natsListener: received order %q", order.OrderUID)
}
