package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	InitDebounce()
	InitThrottle()
}

func Debounce(f func(), d time.Duration) func() {
	var mu sync.Mutex
	var timer *time.Timer

	return func() {
		mu.Lock()
		defer mu.Unlock()

		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(d, f)
	}
}

func InitDebounce() {
	search := Debounce(func() {
		fmt.Println("Searching...")
	}, 500*time.Millisecond)

	fmt.Println("Typing 'go'...")
	search()
	time.Sleep(100 * time.Millisecond)
	search()
	time.Sleep(100 * time.Millisecond)
	search()

	fmt.Println("Typing 'lang'...")
	time.Sleep(400 * time.Millisecond)
	search()

	time.Sleep(1 * time.Second)
}

func Throttle(f func(), d time.Duration) func() {
	var mu sync.Mutex
	lastCalled := time.Now().Add(-d) 

	return func() {
		mu.Lock()
		defer mu.Unlock()

		if time.Since(lastCalled) >= d {
			lastCalled = time.Now()
			go f() 
		}
	}
}

func InitThrottle() {
	click := Throttle(func() {
		fmt.Println("Button clicked!")
	}, 1*time.Second)

	fmt.Println("Clicking button 5 times in 2 seconds...")
	click()
	time.Sleep(200 * time.Millisecond)
	click()
	time.Sleep(200 * time.Millisecond)
	click()
	time.Sleep(200 * time.Millisecond)
	click()
	time.Sleep(200 * time.Millisecond)
	click()

	time.Sleep(1 * time.Second)
	click()

	time.Sleep(2 * time.Second)
}