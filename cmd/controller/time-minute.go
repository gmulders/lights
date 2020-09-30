package main

import (
	"encoding/json"
	"time"

	"github.com/kelvins/sunrisesunset"
	log "github.com/sirupsen/logrus"
)

func parseTimeMinute(data []byte) (interface{}, error) {
	event := timeMinuteEvent{}
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

type timeMinuteEvent struct {
	Time time.Time `json:"time"`
}

func (event timeMinuteEvent) SunSetsTimePlusMinutes(lat float64, lon float64, minutes int) bool {
	t := event.Time
	t = t.Add(time.Duration(minutes) * time.Minute)
	_, offset := t.Zone()
	utcOffset := float64(offset) / float64(3600)

	p := sunrisesunset.Parameters{
		Latitude:  lat,
		Longitude: lon,
		UtcOffset: utcOffset,
		Date:      time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC),
	}

	_, sunset, err := p.GetSunriseSunset()

	if err != nil {
		log.Errorf("Could not determine sunset: %v", err)
		return false
	}

	result := t.Hour() == sunset.Hour() && t.Minute() == sunset.Minute()
	log.Debugf("The result of SunSetsTimePlusMinutes: %v", result)
	return result
}

func (event timeMinuteEvent) TimeEquals(hour int, minute int) bool {
	t := event.Time
	result := t.Hour() == hour && t.Minute() == minute
	log.Debugf("The result of TimeEquals: %v", result)
	return result
}
