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

const databaseURI = "postgresql://postgres:secret@localhost:5432/wildberries_l0"

func main() {
	pg, err := postgres.NewStorage(databaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer logIfError(pg.Close)
	s, err := inmem.NewStorage(inmem.WithPersistentStorage(pg))
	must(err)
	nl, err := nats_listener.New("cluster-L0", "orderServer", "orderServerSub", "orders", s)
	must(err)
	defer logIfError(nl.Close)

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

func logIfError(closeFn func() error) {
	if err := closeFn(); err != nil {
		log.Println(err)
	}
}
