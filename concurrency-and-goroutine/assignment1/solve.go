package main


import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- string) {
	for j := range jobs {
		fmt.Printf("worker %d starting job %d\n", id, j)
		time.Sleep(time.Millisecond * 500)
		fmt.Printf("worker %d finished job %d\n", id, j)
		results <- fmt.Sprintf("Job %d completed by worker %d", j, id)
	}
}

func main() {
	const numJobs = 10
	const numWorkers = 3
	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	close(jobs)

	for a := 1; a <= numJobs; a++ {
		fmt.Println(<-results)
	}
}
