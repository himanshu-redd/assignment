package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(ch chan<- int) {
	for i := 1; i <= 10; i++ {
		fmt.Printf("Producer: sending task %d\n", i)
		ch <- i
		time.Sleep(100 * time.Millisecond)
	}
	close(ch)
}

func worker(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range ch {
		fmt.Printf("Worker %d: processing task %d\n", id, task)
		time.Sleep(300 * time.Millisecond) // simulate work
	}
}

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, ch, &wg)
	}

	go producer(ch)

	wg.Wait()
}
