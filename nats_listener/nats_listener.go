package nats_listener

import "github.com/nats-io/stan.go"

type NATSListener struct {
	sc  *stan.Conn
	sub *stan.Subscription
}

func New(stanCluster, clientID, durableName, subject string) (NATSListener, error) {
	sc, err := stan.Connect(stanCluster, clientID)
	if err != nil {
		return NATSListener{}, err
	}
	nl := NATSListener{
		sc: &sc,
	}
	sub, err := sc.Subscribe(subject, nl.msgHandler, stan.DurableName(durableName))
	if err != nil {
		return NATSListener{}, err
	}
	nl.sub = &sub

	return nl, nil
}

func (nl NATSListener) msgHandler(msg *stan.Msg) {}
