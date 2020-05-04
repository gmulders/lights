package main

import (
	"time"

	"github.com/gmulders/lights"
	nats "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Errorf("Could not connect to NATS: %v", err)
		return
	}
	defer nc.Close()

	for t := range lights.MinuteTicker().C {
		err = nc.Publish("time-minute", []byte(`{"time":"` + t.Format(time.RFC3339) + `"}`))
		if err != nil {
			log.Errorf("Could not publish 'time-minute' event to NATS: %v", err)
			return
		}
	}
}
