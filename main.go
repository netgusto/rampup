package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	// "strconv"
	"time"
)

// 1. Get list of URLS
// 2. 10x in a row, call each URL of the list and get the duration, reachability, status code of each call
// 3. Try n times in case of error for each URL
// 4. Display metrics on stdout

func main() {
	urls := getURLList()
	spew.Dump(urls)
	for k := 0; k < 10; k++ {
		for _, url := range urls {
			fmt.Println(url)

			measure := measureURL(url, 3)
			spew.Dump(measure)
			// fmt.Printf("Status code is %s, duration is %s and also err is %v (at try %d)\n", statusCodeStr, d, err, try)
			panic("")
		}
	}
}

type measureResult struct  {
	statusCode *int
	retries int
	duration time.Duration
	err error
}

func measureURL(url string, maxRetries int) measureResult {
	var statusCode *int
	var retries int
	var duration time.Duration
	var err error
	for retries = 0; retries <= maxRetries; retries++ {
		t1 := time.Now()
		var resp *http.Response
		resp, err = http.Get(url)
		duration = time.Since(t1)
		if err == nil {
			statusCode = &resp.StatusCode
			break
		}
		time.Sleep(1 * time.Second)
	}
	return measureResult{statusCode: statusCode, retries: retries - 1, duration: duration, err: err}
}

func getURLList() []string {
	// FIXME: implement
	return []string{"https://www.algoliaaeraerazerazerazerazerazer.com",
		"https://d85-usw-1.algolia.net/1/isalive",
		"https://d85-usw-2.algolia.net/1/isalive",
		"https://d85-usw-3.algolia.net/1/isalive",
	}
}
