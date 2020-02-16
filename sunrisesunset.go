package lights

import (
	"github.com/kelvins/sunrisesunset"
	"log"
	"time"
)

// IsSunSet returns wether the sun has set at the given time.
func IsSunSet(t time.Time) (bool, error) {
	_, offset := t.Zone()
	utcOffset := float64(offset) / float64(3600)

	// You can use the Parameters structure to set the parameters
	p := sunrisesunset.Parameters{
		Latitude:  52.09083,
		Longitude: 05.12222,
		UtcOffset: utcOffset,
		Date:      time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC),
	}

	_, sunset, err := p.GetSunriseSunset()

	sunset = time.Date(t.Year(), t.Month(), t.Day(), sunset.Hour(), sunset.Minute(), sunset.Second(), sunset.Nanosecond(), time.UTC)

	if err != nil {
		log.Printf("Bla %s", err)
		return true, err
	}

	// log.Printf("Current time: %s, sunset time: %s", t.Format(time.RFC3339), sunset.Format(time.RFC3339))

	return t.After(sunset), nil
}
