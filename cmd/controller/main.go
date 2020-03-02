package main

import (
	"encoding/json"

	"github.com/gmulders/lights"
	nats "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

var (
	sunsetLights = []lights.Light{
		lights.Light{
			ID:  2,
			On:  true,
			Bri: 170,
			Sat: 0,
		},
		lights.Light{
			ID:  5,
			On:  true,
			Bri: 254,
		},
		lights.Light{
			ID:  6,
			On:  true,
		},
	}

	time2030Lights = []lights.Light{
		lights.Light{
			ID: 2,
			On: false,
		},
	}

	time2045Lights = []lights.Light{
		lights.Light{
			ID:  6,
			On:  false,
		},
	}

	time2100Lights = []lights.Light{
		lights.Light{
			ID: 5,
			On: false,
		},
	}
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Could not connect to NATS: %v", err)
	}

	defer nc.Close()

	nc.Subscribe("sunset", func(msg *nats.Msg) {
		log.Info("Received sunset event")
		updateLights(nc, sunsetLights)
	})

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Errorf("Unexpected exception from NATS: %v", err)
	}

	for t := range lights.MinuteTicker().C {
		currentTimeOnly := t.Hour()*100 + t.Minute()
		if currentTimeOnly == 2030 {
			updateLights(nc, time2030Lights)
		} else if currentTimeOnly == 2100 {
			updateLights(nc, time2100Lights)
		}

		nc.Flush()
	}
}

func updateLights(nc *nats.Conn, lights []lights.Light) {
	bytes, err := json.Marshal(lights)

	if err != nil {
		log.Errorf("Could not marshal to JSON %v", err)
		return
	}

	log.Infof("Sending request for 'lights-update' %s", string(bytes))

	nc.Publish("lights-update", bytes)
}
