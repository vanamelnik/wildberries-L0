package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vanamelnik/wildberries-L0/nats_listener"
	"github.com/vanamelnik/wildberries-L0/storage/inmem"
	"github.com/vanamelnik/wildberries-L0/storage/postgres"
)

const databaseURI = "postgresql://postgres:qwe123@localhost:5432/wildberries_l0"

func main() {
	pg, err := postgres.NewStorage(databaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()
	s := inmem.NewStorage()
	nl, err := nats_listener.New("cluster-L0", "orderServer", "orderServerSub", "orders", s)
	must(err)
	defer nl.Close()

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
