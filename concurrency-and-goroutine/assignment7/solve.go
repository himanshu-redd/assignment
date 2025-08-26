package main

import (
	"context"
	"fmt"
	"time"
)

func longRunningJob(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Job stopped:", ctx.Err())
			return
		default:
			fmt.Println("Working...")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go longRunningJob(ctx)

	time.Sleep(3 * time.Second)
	cancel()

	time.Sleep(1 * time.Second)
	fmt.Println("Main finished")
}
