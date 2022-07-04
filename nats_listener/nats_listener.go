package nats_listener

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/nats-io/stan.go"
	"github.com/vanamelnik/wildberries-L0/models"
	"github.com/vanamelnik/wildberries-L0/storage"
)

// NATSListener is used for listening to the NATS streaming server.
// When a message (order) from the given subject is received, it is stored in the storage provided.
type NATSListener struct {
	sc  stan.Conn
	sub stan.Subscription
	s   storage.Storage
}

// New creates a new connection to the nats-streaming-server and registers a callback method that
// processes incoming orders.
func New(stanCluster, clientID, durableName, subject string, s storage.Storage) (NATSListener, error) {
	sc, err := stan.Connect(stanCluster, clientID)
	if err != nil {
		return NATSListener{}, err
	}
	nl := NATSListener{
		sc: sc,
		s:  s,
	}
	sub, err := sc.Subscribe(subject, nl.msgHandler, stan.DurableName(durableName))
	if err != nil {
		return NATSListener{}, err
	}
	nl.sub = sub

	return nl, nil
}

// Close closes nats streaming subscription and connection.
func (nl NATSListener) Close() (retErr error) {
	if err := nl.sub.Close(); err != nil {
		retErr = multierror.Append(retErr, err)
	}
	if err := nl.sc.Close(); err != nil {
		retErr = multierror.Append(retErr, err)
	}
	return
}

// msgHandler is a callback function that sends all incoming orders to the storage.
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
	if err := nl.s.Store(order.OrderUID, string(msg.Data)); err != nil {
		log.Printf("natsListener: ERR: order rejected: %s", err)
	}
	log.Printf("natsListener: order %q received and stored", order.OrderUID)
}
