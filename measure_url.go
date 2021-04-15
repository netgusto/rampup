package main

import (
	"time"
)

type measureResult struct {
	statusCode *int
	tries      int
	duration   time.Duration
	err        error
}

func measureURL(url string, maxTries int, statusGetter URLStatusGetter) measureResult {
	var statusCode *int
	var duration time.Duration
	var err error

	tries := 0
	for i := 0; i < maxTries; i++ {
		tries++

		t1 := time.Now()
		statusCode, err = statusGetter.GetStatus(url)
		duration = time.Since(t1)
		if err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return measureResult{
		statusCode: statusCode,
		tries:      tries,
		duration:   duration,
		err:        err,
	}
}
