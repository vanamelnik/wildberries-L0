package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vanamelnik/wildberries-L0/nats_listener"
	"github.com/vanamelnik/wildberries-L0/server"
	"github.com/vanamelnik/wildberries-L0/storage/inmem"
	"github.com/vanamelnik/wildberries-L0/storage/postgres"
)

const (
	databaseURI = "postgresql://postgres:secret@localhost:5432/wildberries_l0"
	addr        = ":8080"
)

func main() {
	pg, err := postgres.NewStorage(databaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer logIfError(pg.Close)
	s, err := inmem.NewCache(inmem.WithPersistentStorage(pg))
	must(err)
	nl, err := nats_listener.New("cluster-L0", "orderServer", "orderServerSub", "orders", s)
	must(err)
	defer logIfError(nl.Close)

	log.Println("NATS Listener started")

	server, err := server.New(addr, s)
	must(err)
	go logIfError(server.ListenAndServe)
	log.Printf("HTTP server is listening at %s", addr)
	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
	}()

	sigint := make(chan os.Signal, 1)
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
