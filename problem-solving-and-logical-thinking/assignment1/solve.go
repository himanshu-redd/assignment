package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	const (
		numProducers = 3
		numConsumers = 5
		bufferSize   = 10
		maxItems     = 20
	)

	dataChan := make(chan int, bufferSize)

	var wg sync.WaitGroup

	wg.Add(numProducers)
	for i := 1; i <= numProducers; i++ {
		go producer(i, dataChan, &wg, maxItems/numProducers)
	}

	wg.Add(numConsumers)
	for i := 1; i <= numConsumers; i++ {
		go consumer(i, dataChan, &wg)
	}

	go func() {
		wg.Wait()
		close(dataChan)
	}()

	fmt.Println("Program started. Producers and consumers are running.")
	fmt.Println("Program finished.")
}

func producer(id int, dataChan chan<- int, wg *sync.WaitGroup, count int) {
	defer wg.Done()

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		item := rand.Intn(100) + 1
		fmt.Printf("Producer %d: sending item %d\n", id, item)
		dataChan <- item
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
	fmt.Printf("Producer %d: finished producing\n", id)
}

func consumer(id int, dataChan <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range dataChan {
		fmt.Printf("Consumer %d: received item %d\n", id, item)
		time.Sleep(time.Duration(rand.Intn(150)) * time.Millisecond)
	}
	fmt.Printf("Consumer %d: finished consuming\n", id)
}