package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	nats "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

var (
	actionMap map[string][]eventAction
	nc *nats.Conn
)

func main() {
	if err := readConfig(); err != nil {
		log.Errorf("Could not read config: %v", err)
		return
	}

	var err error

	// Connect to NATS
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Could not connect to NATS: %v", err)
	}
	defer nc.Close()

	nc.Subscribe("time-minute", processEvent(parseTimeMinute))

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Errorf("Unexpected exception from NATS: %v", err)
	}

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
}

func readConfig() error {
	viper.SetConfigName("config")
	viper.BindEnv("hue.username")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	buildActionMap()

	viper.WatchConfig()
	viper.OnConfigChange(func(event fsnotify.Event) {
		if event.Op != fsnotify.Write {
			return
		}
		buildActionMap()
	})

	return nil
}

func buildActionMap() error {
	actions := make([]eventAction, 0)
	if err := viper.UnmarshalKey("actions", &actions); err != nil {
		return err
	}

	actionMap = map[string][]eventAction{}

	for i := 0; i < len(actions); i++ {
		action := actions[i]
		actionMap[action.Subject] = append(actionMap[action.Subject], action)
	}

	log.Infof("Config: %v", actionMap)

	return nil
}
