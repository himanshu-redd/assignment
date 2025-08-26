package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const queueSize = 5

var wg sync.WaitGroup

func producer(id int, dataChan chan<- int, waitChan chan <- bool) {

	rand.Seed(time.Now().UnixNano() + int64(id))

	for i := 0; i < 5; i++ {
		num := rand.Intn(90) + 10
		
		fmt.Printf("Producer %d: sending value %d\n", id, num)
		
		dataChan <- num
		
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	}

	close(dataChan)
	waitChan <- true
}

func consumer(id int, dataChan <-chan int) {

	for num := range dataChan {
		fmt.Printf("Consumer %d: received value %d\n", id, num)
		
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	}
	fmt.Printf("Consumer %d: channel closed, exiting.\n", id)
}

func main() {
	dataChan := make(chan int, queueSize)
	waitChan := make(chan bool)

	numProducers := 3
	numConsumers := 2


	for i := 1; i <= numProducers; i++ {
		go producer(i, dataChan, waitChan)
	}

	for i := 1; i <= numConsumers; i++ {
		go consumer(i, dataChan)
	}

	<- waitChan

	fmt.Println("All producers finished. Waiting for consumers to finish...")
	fmt.Println("Program finished.")
}
