package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(id int, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		item := id*10 + i
		fmt.Printf("Producer %d: producing item %d\n", id, item)
		ch <- item
		time.Sleep(100 * time.Millisecond)
	}
}

func consumer(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range ch {
		fmt.Printf("Consumer %d: consuming item %d\n", id, item)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Printf("Consumer %d: channel closed, exiting\n", id)
}

func main() {
	const numProducers = 2
	const numConsumers = 3

	ch := make(chan int, 10)
	var wgProducer sync.WaitGroup
	var wgConsumer sync.WaitGroup

	wgConsumer.Add(numConsumers)
	for i := 1; i <= numConsumers; i++ {
		go consumer(i, ch, &wgConsumer)
	}

	wgProducer.Add(numProducers)
	for i := 1; i <= numProducers; i++ {
		go producer(i, ch, &wgProducer)
	}

	wgProducer.Wait()
	close(ch)
	wgConsumer.Wait()

	fmt.Println("All producers and consumers have finished.")
}
