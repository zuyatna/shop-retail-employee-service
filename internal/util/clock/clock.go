package clock

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (rc RealClock) Now() time.Time {
	return time.Now()
}
