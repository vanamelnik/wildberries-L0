package main

// orderpub is a publisher of json orders to nats-streaming-server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

const (
	clusterName = "cluster-L0"
	clientID    = "orderPub"
	subject     = "orders"
)

func main() {
	sc, err := stan.Connect(clusterName, clientID)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()
	log.Println("Successfully connected to nats-streaming")

	var order string
	if len(os.Args) < 2 {
		order = readFromConsole()
	} else {
		order = readFromFile()
	}
	if err := sc.Publish(subject, []byte(order)); err != nil {
		log.Fatal(err)
	}
	log.Println("Order sent")
}

func readFromConsole() string {
	fmt.Println("Type the order manually:")
	r := bufio.NewReader(os.Stdin)
	res, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func readFromFile() string {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	res, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order successfully read from the file %s\n", os.Args[1])
	return string(res)
}
