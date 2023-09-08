package shared_lib

import "time"

func Now() time.Time {
	now := time.Now()
	tz, _ := time.LoadLocation("Europe/London")
	now.In(tz)
	return now
}
