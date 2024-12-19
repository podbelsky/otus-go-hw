package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(0)

	timeout := flag.Duration("timeout", 10*time.Second, "connect timeout")

	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatal("Host and port arguments must be provided")
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	tc := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := tc.Connect(); err != nil {
		log.Fatalf("Failed to connect: %s", err)
	}

	log.Println("...Connected to", address)

	defer func() {
		if err := tc.Close(); err != nil {
			log.Println(err)
		}
	}()

	notifyContext, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ctx, cancelFn := context.WithCancel(notifyContext)

	go send(tc, cancelFn)
	go receive(tc, cancelFn)

	<-ctx.Done()
}

func send(client TelnetClient, cancel context.CancelFunc) {
	if err := client.Send(); err != nil {
		log.Println(err)
	}

	cancel()
}

func receive(client TelnetClient, cancel context.CancelFunc) {
	if err := client.Receive(); err != nil {
		log.Println(err)
	}

	cancel()
}
