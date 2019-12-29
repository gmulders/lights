package main

import (
	"log"
	"runtime"

	nats "github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	nc.Subscribe("sunset", func(msg *nats.Msg) {
		printMsg(msg)
	})

	nc.Subscribe("lights-update", func(msg *nats.Msg) {
		nc.Publish("lights-status", []byte("TODO"))
	})

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags)

	runtime.Goexit()
}

func printMsg(m *nats.Msg) {
	log.Printf("Received on [%s]: '%s'", m.Subject, string(m.Data))
}
