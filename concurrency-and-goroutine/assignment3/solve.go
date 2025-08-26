package main

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	tokens chan struct{}
}

func NewRateLimiter(rate int, burst int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, burst),
	}

	refillTicker := time.NewTicker(time.Second / time.Duration(rate))

	go func() {
		for range refillTicker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()

	for i := 0; i < burst; i++ {
		rl.tokens <- struct{}{}
	}

	return rl
}

func (rl *RateLimiter) Allow() {
	<-rl.tokens
}

func main() {
	rate := 5
	burst := 5
	limiter := NewRateLimiter(rate, burst)

	var wg sync.WaitGroup
	const totalRequests = 20

	fmt.Printf("Starting a rapid burst of %d requests...\n", totalRequests)

	for i := 1; i <= totalRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			limiter.Allow()
			fmt.Printf("Processing request %d at %s\n", id, time.Now().Format("15:04:05.000"))
		}(i)
	}

	wg.Wait()
	fmt.Println("All requests processed.")
}
