package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func doRequestWithTimeout(url string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

func main() {
	slowURL := "https://httpbin.org/delay/5"
	fastURL := "https://httpbin.org/get"

	fmt.Printf("Attempting request to %s with 3-second timeout...\n", slowURL)
	body, err := doRequestWithTimeout(slowURL, 3*time.Second)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Response: %s\n", body)
	}

	fmt.Println("--------------------")

	fmt.Printf("Attempting request to %s with 3-second timeout...\n", fastURL)
	body, err = doRequestWithTimeout(fastURL, 3*time.Second)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Response: %s\n", body)
	}
}
