package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var urls = []string{
	"http://google.ru",
	"wrong-",
	"unexisting.ru",
	"http://ya.ru",
	"http://vk.com",
	"http://vk.com",
	"http://vk.com",
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- int) {
	defer wg.Done()
	for {
		select {
		case u, ok := <-jobs:
			if !ok {
				fmt.Println("stopping worker!")
				return
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
			if err != nil {
				continue
			}

			httpClient := &http.Client{
				Timeout: time.Second * 5,
			}

			r, err := httpClient.Do(req)
			if err != nil {
				continue
			}
			defer r.Body.Close()

			results <- r.StatusCode
		case <-ctx.Done():
			fmt.Println("canceling worker!")
			return
		}
	}
}

func Start(urls []string) {
	var count int
	var wg sync.WaitGroup

	n := 100

	jobs := make(chan string, n)
	results := make(chan int, n)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(ctx, &wg, jobs, results)

	}

	go func() {
		for _, u := range urls {
			select {
			case <-ctx.Done():
				fmt.Println("stop writing url")
				return
			case jobs <- u:
			}
		}
		close(jobs)
	}()

	go func() {
		for r := range results {
			if r == 200 {
				count++
				if count == 2 {
					fmt.Println("canceling")
					cancel()
					break
				}
			}
		}
	}()
	wg.Wait()
	close(results)
	fmt.Println(count)
}

func main() {
	n := 189695511
	u := make([]string, n)
	for i := 0; i < n; i++ {
		u[i] = urls[i%7]
	}
	Start(u)
}
