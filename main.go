package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

func processUrl(url string) int {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	bodyString := string(body)

	return strings.Count(bodyString, "Go")
}

func main() {

	const GOROUTINES_COUNT = 5
	const URLS_COUNT = 5

	urls := make([]string, URLS_COUNT)
	results := make([]int, URLS_COUNT)
	total := 0

	wg := &sync.WaitGroup{}
	semaphore := make(chan bool, GOROUTINES_COUNT)

	for i := 0; i < len(urls); i++ {
		urls[i] = "https://golang.org"
	}

	wg.Add(len(urls))

	for i := 0; i < len(urls); i++ {
		semaphore <- true

		go func(i int, url string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			results[i] = processUrl(url)
		}(i, urls[i])
	}

	wg.Wait()

	for i := 0; i < len(urls); i++ {
		fmt.Printf("Count for %v: %v \n", urls[i], results[i])
		total += results[i]
	}

	fmt.Printf("Total: %v \n", total)
}
