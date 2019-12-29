package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/amimof/huego"
	"github.com/gmulders/lights"
	nats "github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"log"
	"runtime"
	"strings"
)

func main() {

	config, err := readConfig()

	if err != nil {
		log.Fatal(err)
	}

	bridge, err := connectToHueBridge(config)

	if err != nil {
		log.Fatal(err)
	}

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	nc.Subscribe("sunset", func(msg *nats.Msg) {
		bridge.SetLightState(2, huego.State{On: true, Bri: 170, Sat: 0})
	})

	nc.Subscribe("lights-update", func(msg *nats.Msg) {
		bytes, err := getAllLightsAsJSON(bridge)
		if err != nil {
			log.Fatal(err)
		}
		nc.Publish("lights-status", bytes)
	})

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags)

	runtime.Goexit()
}

type config struct {
	username string
	ip       string
}

func readConfig() (*config, error) {
	viper.SetConfigName("config")
	viper.BindEnv("hue.username")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %s", err)
	}

	if !viper.IsSet("hue.username") {
		return nil, errors.New("Missing hue username")
	}

	username := viper.GetString("hue.username")
	ip := viper.GetString("hue.ip")

	return &config{username: username, ip: ip}, nil
}

func connectToHueBridge(config *config) (*huego.Bridge, error) {
	if config.ip != "" {
		log.Print("Connect using ip")
		return huego.New(config.ip, config.username), nil
	}

	log.Print("Connect using discovery")
	bridge, err := huego.Discover()
	if err != nil {
		return nil, errors.New("Can't find the hue bridge")
	}
	return bridge.Login(config.username), nil
}

func getAllLightsAsJSON(bridge *huego.Bridge) ([]byte, error) {
	oldLights, err := bridge.GetLights()

	if err != nil {
		return nil, fmt.Errorf("could not get lights %s", err)
	}

	newLights := make([]lights.Light, len(oldLights))
	for index, light := range oldLights {
		newLights[index] = lights.Light{
			ID:   light.ID,
			Name: light.Name,
			Hue:  light.State.Hue,
			Sat:  light.State.Sat,
			Bri:  light.State.Bri,
		}
	}

	b, err := json.Marshal(newLights)

	if err != nil {
		return nil, fmt.Errorf("could not marshal to JSON %s", err)
	}

	return b, nil
}
