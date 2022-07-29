package GameTime

import "time"

var startTime time.Time

func init() {
	startTime = time.Now()
}

func getTime() float64 {
	dur := time.Now().Sub(startTime)
	return dur.Seconds()
}
