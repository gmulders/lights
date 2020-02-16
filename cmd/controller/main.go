package main

import (
	"encoding/json"
	"fmt"
	"github.com/gmulders/lights"
	nats "github.com/nats-io/nats.go"
	"log"
	"runtime"
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
			Bri: 136,
		},
	}

	time2030Lights = []lights.Light{
		lights.Light{
			ID: 2,
			On: false,
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
		log.Fatal(err)
	}

	nc.Subscribe("sunset", func(msg *nats.Msg) {
		updateLights(nc, sunsetLights)
	})

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags)

	for t := range lights.MinuteTicker().C {
		currentTimeOnly := t.Hour()*100 + t.Minute()
		if currentTimeOnly == 2030 {
			updateLights(nc, time2030Lights)
		} else if currentTimeOnly == 2100 {
			updateLights(nc, time2100Lights)
		}
	}

	runtime.Goexit()
}

func updateLights(nc *nats.Conn, lights []lights.Light) {
	bytes, err := json.Marshal(lights)

	if err != nil {
		log.Print(fmt.Errorf("could not marshal to JSON %s", err))
	}

	nc.Publish("lights-update", bytes)
}
