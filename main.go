package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"strconv"
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

			t1 := time.Now()
			resp, err := http.Get(url)
			d := time.Since(t1)
			statusCodeStr := "ERROR"
			if err == nil {
				statusCodeStr = strconv.Itoa(resp.StatusCode)
			}
			fmt.Printf("Status code is %s, duration is %s and also err is %v\n", statusCodeStr, d, err)
		}
	}
}

func getURLList() []string {
	// FIXME: implement
	return []string{"https://www.algoliaaeraerazerazerazerazerazer.com",
		"https://d85-usw-1.algolia.net/1/isalive",
		"https://d85-usw-2.algolia.net/1/isalive",
		"https://d85-usw-3.algolia.net/1/isalive",
	}
}
