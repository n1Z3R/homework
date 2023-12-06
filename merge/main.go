package main

import (
	"fmt"
	"sync"
)

func merge[T any](chans ...<-chan T) <-chan T {
	res := make(chan T)
	wg := sync.WaitGroup{}

	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch <-chan T) {
			defer wg.Done()
			for r := range ch {
				res <- r
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}
func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		ch1 <- 1
		ch1 <- 2
		close(ch1)
	}()
	go func() {
		ch2 <- 1
		ch2 <- 2
		close(ch2)
	}()

	for v := range merge(ch1, ch2) {
		fmt.Println(v) // 1,2,3,4 в любом порядке. Цикл корректно завершается.
	}
}
