package main

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	tokens       float64
	capacity     float64
	refillRate   float64
	lastFillTime time.Time
	mu           sync.Mutex
}

func NewRateLimiter(capacity, refillRate float64) *RateLimiter {
	return &RateLimiter{
		capacity:     capacity,
		refillRate:   refillRate,
		tokens:       capacity,
		lastFillTime: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastFillTime).Seconds()

	tokensToAdd := elapsed * rl.refillRate

	rl.tokens = min(rl.tokens+tokensToAdd, rl.capacity)

	rl.lastFillTime = now

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}

	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	limiter := NewRateLimiter(10, 2)

	fmt.Println("Simulating 20 requests with a rate limiter (Capacity: 10, Refill: 2/sec)...")

	for i := 1; i <= 20; i++ {
		if limiter.Allow() {
			fmt.Printf("Request #%d: ALLOWED\n", i)
		} else {
			fmt.Printf("Request #%d: DENIED\n", i)
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nWaiting for 3 seconds to let the bucket refill...")
	time.Sleep(3 * time.Second)

	fmt.Println("\nSimulating more requests after the wait...")

	for i := 21; i <= 25; i++ {
		if limiter.Allow() {
			fmt.Printf("Request #%d: ALLOWED\n", i)
		} else {
			fmt.Printf("Request #%d: DENIED\n", i)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
