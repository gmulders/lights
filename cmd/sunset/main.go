package main

import (
	"github.com/gmulders/lights"
	nats "github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	isSunset, err := lights.IsSunSet(time.Now())

	if err != nil {
		log.Fatal(err)
		return
	}

	for t := range lights.MinuteTicker().C {
		isSunsetNew, err := lights.IsSunSet(time.Now())

		if err != nil {
			log.Fatal(err)
			break
		}

		if isSunsetNew == true && isSunset == false {
			sendSunsetEvent(t)
		}

		isSunset = isSunsetNew
	}

}

func sendSunsetEvent(t time.Time) {
	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer nc.Close()
	err = nc.Publish("sunset", []byte(t.Format(time.RFC3339)))
	if err != nil {
		log.Fatal(err)
		return
	}
}
