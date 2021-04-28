package main

import (
	"context"
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigs
		fmt.Println("Signal caught, terminating.")
		cancel()
	}()

	measures := make(chan measureResult)
	go consume(ctx, measures)

	loop(ctx, urls, statusGetter, measures, time.Second*10)

	fmt.Println("Terminated.")
}

func consume(ctx context.Context, measures <-chan measureResult) {
	for {
		select {
		case <-ctx.Done():
			return
		case measure := <-measures:
			fmt.Printf("URL: %s, status: %d, tries: %d, duration: %s, err: %v\n", measure.url, *measure.statusCode, measure.tries, measure.duration, measure.err)
		}
	}
}

func worker(ctx context.Context, workerID int, jobs <-chan string, statusGetter URLStatusGetter, measures chan<- measureResult) {
	for url := range jobs {
		select {
		case <-ctx.Done():
			fmt.Printf("Exiting worker %d\n", workerID)
			return
		case measures <- measureURL(url, 3, statusGetter):
		}
	}
}

func loop(ctx context.Context, urls []string, statusGetter URLStatusGetter, measures chan<- measureResult, minRoundDuration time.Duration) {

	jobs := make(chan string)
	for k := 0; k < 2; k++ {
		go worker(ctx, k, jobs, statusGetter, measures)
	}

	round := 0

	for {
		start := time.Now()

		round += 1
		fmt.Printf("Starting round %d\n", round)

		for _, url := range urls {
			select {
			case <-ctx.Done():
				fmt.Println("Exiting!")
				return
			case jobs <- url:
			}
		}

		thisRoundDuration := time.Since(start)
		cooldown := time.Second * 0

		if thisRoundDuration < minRoundDuration {
			cooldown = minRoundDuration - thisRoundDuration
			fmt.Printf("Looping to fast! Waiting %s\n", cooldown)
		}

		select {
		case <-ctx.Done():
			fmt.Println("Exiting!")
			return
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
