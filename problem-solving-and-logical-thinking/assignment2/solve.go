package main

import (
	"fmt"
	"time"
)

type Job struct {
	ID int
	A  int
	B  int
}

type Result struct {
	JobID int
	Sum   int
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, job.ID)
		time.Sleep(time.Millisecond * 100)
		sum := job.A + job.B
		result := Result{JobID: job.ID, Sum: sum}
		results <- result
		fmt.Printf("Worker %d finished job %d\n", id, job.ID)
	}
}

func main() {
	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, A: j, B: j * 2}
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		result := <-results
		fmt.Printf("Job %d completed with result %d\n", result.JobID, result.Sum)
	}
}