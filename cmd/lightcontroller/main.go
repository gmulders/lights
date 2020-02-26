package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/amimof/huego"
	"github.com/gmulders/lights"
	nats "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	config, err := readConfig()

	if err != nil {
		log.Fatalf("Could not read config: %v", err)
	}

	bridge, err := connectToHueBridge(config)

	if err != nil {
		log.Fatalf("Could not connect to bridge: %v", err)
	}

	go updateLightsStatus(bridge)

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Could not connect to NATS: %v", err)
	}
	defer nc.Close()

	nc.Subscribe("lights-update", func(msg *nats.Msg) {
		lights := make([]lights.Light, 0)
		if err := json.Unmarshal(msg.Data, &lights); err != nil {
			log.Errorf("Unable to unmarshal the lights JSON: %v", err)
			return
		}

		for _, light := range lights {
			updateLight(bridge, light)
		}
	})

	nc.Subscribe("lights-status-query", func(msg *nats.Msg) {
		bytes, err := getAllLightsAsJSON(bridge)
		if err != nil {
			log.Error(err)
			return
		}
		nc.Publish("lights-status-response", bytes)
	})

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Error(err)
	}

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
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

var currentLights []huego.Light

func getAllLightsAsJSON(bridge *huego.Bridge) ([]byte, error) {
	bridgeLights, err := bridge.GetLights()

	if err != nil {
		return nil, fmt.Errorf("could not get lights %s", err)
	}

	currentLights = bridgeLights

	newLights := make([]lights.Light, len(currentLights))
	for index, light := range currentLights {
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

func updateLightsStatus(bridge *huego.Bridge) {
	ticker := time.NewTicker(5 * time.Minute)

	for ; true; <-ticker.C {
		bridgeLights, err := bridge.GetLights()
		if err != nil {
			log.Print(fmt.Errorf("could not get lights %s", err))
		}
		currentLights = bridgeLights
	}
}

func findLightIDByName(name string) (int, error) {
	for _, light := range currentLights {
		if light.Name == name {
			return light.ID, nil
		}
	}
	return 0, fmt.Errorf("light with name '%s' could not be found", name)
}

func updateLight(bridge *huego.Bridge, light lights.Light) {
	if light.ID == 0 {
		id, err := findLightIDByName(light.Name)
		if err != nil {
			log.Error(err)
			return
		}
		light.ID = id
	}

	state := huego.State{
		On:  light.On,
		Bri: light.Bri,
		Hue: light.Hue,
		Sat: light.Sat,
	}

	_, err := bridge.SetLightState(light.ID, state)

	if err != nil {
		log.Errorf("Could not set the state of the light %d: %v", light.ID, err)
	}
}
