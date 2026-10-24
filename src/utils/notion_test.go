package utils

import (
	"context"
	"fmt"
	notionModels "notion_saver/src/models/notion"
	"sync"
	"testing"
	"time"
)

func TestSearchPages(t *testing.T) {
	env := LoadEnv("../../.env")
	ctx := context.Background()

	server := Server{
		Env:    env,
		Logger: NewLogger(ctx),
	}

	notion := NewNotion(&server)
	pages, err := notion.SearchPages(FilterTypePage, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if len(pages.Results) == 0 {
		t.Error("No pages found")
		return
	}
}

func TestReachRateLimits(t *testing.T) {
	env := LoadEnv("../../.env")
	ctx := context.Background()

	server := Server{
		Env:    env,
		Logger: NewLogger(ctx),
	}

	notion := NewNotion(&server)

	deadline := time.Now().Add(10 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	numParallelRequests := 10
	type LimitResult struct {
		pages *notionModels.NotionPages
		err   error
	}
	results := make(chan LimitResult, numParallelRequests)

	notionApiCall := func(i int, ctx context.Context) {
		fmt.Println(fmt.Sprintf("Started routine %v", i))

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			pages, err := notion.SearchPages(FilterTypePage, nil)

			select {
			case results <- LimitResult{pages: pages, err: err}:
			case <-ctx.Done():
				return
			}

		}
	}

	var wg sync.WaitGroup
	for i := 0; i < numParallelRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			notionApiCall(i, ctx)
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result.err)
	}
}
