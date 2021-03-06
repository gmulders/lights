package lights

import (
	"time"
)

// MinuteTicker is a simple Ticker that ticks every minute
func MinuteTicker() *time.Ticker {
	c := make(chan time.Time, 1)
	t := &time.Ticker{C: c}
	go func() {
		for {
			n := time.Now()
			if n.Second() == 0 {
				c <- n
			}
			time.Sleep(time.Second)
		}
	}()
	return t
}
