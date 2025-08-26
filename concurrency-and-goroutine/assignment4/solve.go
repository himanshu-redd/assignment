package main


import (
	"fmt"
	"time"
)

func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)

	for _, ch := range channels {
		go func(c <-chan int) {
			for val := range c {
				out <- val
			}
		}(ch)
	}

	return out
}

func generator(start, step int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := start; ; i += step {
			ch <- i
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return ch
}

func main() {
	ch1 := generator(0, 2) 
	ch2 := generator(1, 2) 

	merged := fanIn(ch1, ch2)

	for i := 0; i < 10; i++ {
		fmt.Println(<-merged)
	}
}
