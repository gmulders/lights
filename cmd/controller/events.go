package main

import (
	"encoding/json"
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	nats "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type eventAction struct {
	Subject     string      `json:"name"`
	Trigger     bool        `json:"trigger"`
	TriggerExpr string      `json:"triggerExpr"`
	TriggerProg *vm.Program `json:"-"`
	ResultEvent string      `json:"resultEvent"`
	Action      interface{} `json:"action,omitempty"`
}

func triggersEvent(ep *eventAction, event interface{}) bool {
	if ep.Trigger {
		return true
	}
	if ep.TriggerExpr == "" {
		return false
	}
	if ep.TriggerProg == nil {
		program, err := expr.Compile(ep.TriggerExpr, expr.Env(event))
		if err != nil {
			log.Errorf("Could not compile expression: %v", err)
			return false
		}
		ep.TriggerProg = program
	}

	output, err := expr.Run(ep.TriggerProg, event)
	if err != nil {
		log.Errorf("Could not evaluate expression: %v", err)
		return false
	}

	return output.(bool)
}

func processEvent(parser func (data []byte) (interface{}, error)) func (msg *nats.Msg) {
	return func (msg *nats.Msg) {
		log.Info("Received event: " + msg.Subject)

		value, err := parser(msg.Data)
		if err != nil {
			log.Errorf("Could not parse message data: %v", err)
		}

		eventActions := actionMap[msg.Subject]

		for i := 0; i < len(eventActions); i++ {
			eventAction := eventActions[i]
			trigger := triggersEvent(&eventAction, value)
			log.Infof("%v", trigger)
			if trigger {
				executeEventAction(&eventAction)
			}
		}
	}
}

func executeEventAction(action *eventAction) {
	bytes, err := json.Marshal(convertInterface(action.Action))

	if err != nil {
		log.Errorf("Could not marshal to JSON %v", err)
		return
	}

	log.Infof("Sending request for '%s' %s", action.ResultEvent, string(bytes))

	if err := nc.Publish(action.ResultEvent, bytes); err != nil {
		log.Errorf("Could not send event: %v", err)
	}
}

func convertInterface(i interface{}) interface{} {
	switch v := i.(type) {
	case map[interface{}]interface{}:
		return convertMap(v)
	default:
		return i
	}
}

func convertMap(m map[interface{}]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			res[fmt.Sprint(k)] = convertMap(v2)
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}
