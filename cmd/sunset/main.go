package main

import (
	"github.com/gmulders/lights"
	nats "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	isSunset, err := lights.IsSunSet(time.Now().UTC())

	if err != nil {
		log.Errorf("Could not determine isSunset, assuming the default (true): %v", err)
	}

	for t := range lights.MinuteTicker().C {
		isSunsetNew, err := lights.IsSunSet(t.UTC())
		if err != nil {
			log.Errorf("Could not determine isSunset: %v", err)
		}

		if isSunsetNew == true && isSunset == false {
			sendSunsetEvent(t)
		}

		isSunset = isSunsetNew
	}
}

func sendSunsetEvent(t time.Time) {
	log.Info("Sending sunset event")

	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Errorf("Could not connect to NATS: %v", err)
		return
	}
	defer nc.Close()
	err = nc.Publish("sunset", []byte(t.Format(time.RFC3339)))
	if err != nil {
		log.Errorf("Could not publish 'sunset' event to NATS: %v", err)
		return
	}
}
