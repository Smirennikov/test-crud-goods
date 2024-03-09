package utils

import "time"

func Debounce(interval time.Duration, wait chan struct{}, cb func()) {
	timer := time.NewTimer(interval)
	for {
		select {
		case <-wait:
			timer.Reset(interval)
		case <-timer.C:
			cb()
		}
	}
}
