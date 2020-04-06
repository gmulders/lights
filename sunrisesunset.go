package lights

import (
	"github.com/kelvins/sunrisesunset"
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

	if err != nil {
		return true, err
	}

	sunset = time.Date(t.Year(), t.Month(), t.Day(), sunset.Hour(), sunset.Minute(), sunset.Second(), sunset.Nanosecond(), time.UTC)

	return t.After(sunset), nil
}
