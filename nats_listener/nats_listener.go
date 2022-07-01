package nats_listener

import (
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/nats-io/stan.go"
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
	log.Printf("natsListener: received msg:\n"+
		"\tSequence: %v\n"+
		"\tSubject: %v\n"+
		"\tReply: %v\n"+
		"\tData: %v\n"+
		"\tTimestamp: %v\n"+
		"\tRedelivered: %v\n"+
		"\tRedeliveryCount: %v\n"+
		"\tCRC32: %v\n\n",
		msg.Sequence,
		msg.Subject,
		msg.Reply,
		string(msg.Data),
		time.Unix(0, msg.Timestamp),
		msg.Redelivered,
		msg.RedeliveryCount,
		msg.CRC32)
}
