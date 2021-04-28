package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
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
			spew.Dump(measure)
		}
	}()

	for k := 0; k < 10; k++ {
		for _, url := range urls {
			go func(url string) {
				measures <- measureURL(url, 3, statusGetter)
			}(url)
		}
	}

	time.Sleep(time.Second * 10)
}

func getURLList() []string {
	// FIXME: implement
	return []string{
		"https://www.algoliaaeraerazerazerazerazerazer.com",
		"https://d85-usw-1.algolia.net/1/isalive",
		"https://d85-usw-2.algolia.net/1/isalive",
		"https://d85-usw-3.algolia.net/1/isalive",
	}
}
