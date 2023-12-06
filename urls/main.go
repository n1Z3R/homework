package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

var (
	counter = 0
	mu      sync.Mutex
)

func main() {
	urls := []string{
		"http://google.ru",
		"wrong-",
		"unexisting.ru",
		"http://ya.ru",
		"http://vk.com",
		"http://vk.com",
		"http://vk.com",
	}

	wg := sync.WaitGroup{}

	// Контекст для отмены запросов
	ctx, cancel := context.WithCancel(context.Background())
	// Канал для уведомления о надобности увеличения счетчика удачных запросов
	ch := make(chan struct{})

	for _, url := range urls {
		wg.Add(1)
		go func(wg *sync.WaitGroup, url string) {
			defer wg.Done()

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return
			}
			req = req.WithContext(ctx)

			client := http.Client{}

			r, err := client.Do(req)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					fmt.Println("context canceled!")
				}
				return
			}

			if r.StatusCode == http.StatusOK {
				fmt.Printf("URL: %s StatusCode: %d\n", url, r.StatusCode)
				ch <- struct{}{}
			}
		}(&wg, url)
	}
	go func() {
		for range ch {
			counter++
			if counter == 2 {
				cancel()
			}
		}
	}()
	wg.Wait()

	fmt.Println("goroutines finished.")
	fmt.Printf("counter = %d\n", counter)
}
