package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	// "strconv"
)

// 1. Get list of URLS
// 2. 10x in a row, call each URL of the list and get the duration, reachability, status code of each call
// 3. Try n times in case of error for each URL
// 4. Display metrics on stdout

func main() {
	urls := getURLList()

	statusGetter := URLStatusGetterReal{}

	measures := make(chan measureResult)

	go func() {
		for measure := range measures {
			fmt.Printf("URL: %s, status: %d, tries: %d, duration: %s, err: %v\n", measure.url, *measure.statusCode, measure.tries, measure.duration, measure.err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	minRoundDuration := time.Second * 10
	cooldown := time.Second * 0

loop:
	for {
		start := time.Now()

		for _, url := range urls {
			select {
			case <-sigs:
				fmt.Println("Exiting!")
				break loop
			default:
				go func(url string) {
					measures <- measureURL(url, 3, statusGetter)
				}(url)
			}
		}

		thisRoundDuration := time.Since(start)
		if thisRoundDuration < minRoundDuration {
			cooldown = minRoundDuration - thisRoundDuration
			fmt.Printf("Looping to fast! Waiting %s\n", cooldown)
		}

		select {
		case <-sigs:
			fmt.Println("Exiting!")
			break loop
		case <-time.After(cooldown):
		}
	}
}

func getURLList() []string {
	// FIXME: implement
	return []string{
		"https://d85-usw-1.algolia.net/1/isalive",
		"https://d85-usw-2.algolia.net/1/isalive",
		"https://d85-usw-3.algolia.net/1/isalive",
		"https://c1-test-1.algolia.net/1/isalive",
		"https://c1-test-2.algolia.net/1/isalive",
		"https://c1-test-3.algolia.net/1/isalive",
		"https://c1-jp-1.algolia.net/1/isalive",
		"https://c1-jp-2.algolia.net/1/isalive",
		"https://c1-jp-3.algolia.net/1/isalive",
		"https://c3-de-1.algolia.net/1/isalive",
		"https://c3-de-2.algolia.net/1/isalive",
		"https://c3-de-3.algolia.net/1/isalive",
	}
}
