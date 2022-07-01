package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vanamelnik/wildberries-L0/nats_listener"
)

func main() {
	nl, err := nats_listener.New("cluster-L0", "orderServer", "orderServerSub", "orders")
	must(err)
	defer func() {
		if err := nl.Close(); err != nil {
			log.Println(err)
		}
		log.Println("App stopped")
	}()

	log.Println("NATS Listener started")

	sigint := make(chan os.Signal)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint
	log.Println("Shutting down...")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
