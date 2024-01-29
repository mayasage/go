package main

import (
	"fmt"
	"sync"
	"time"
)

var sharedLock sync.Mutex
var wg sync.WaitGroup
var runtime = 1 * time.Second

func greedyWorker() {
	defer wg.Done()

	var count int

	for begin := time.Now(); time.Since(begin) <= runtime; {
		sharedLock.Lock()
		time.Sleep(3 * time.Nanosecond)
		sharedLock.Unlock()
		count += 1
	}

	fmt.Printf("Greedy Boy died %d times\n", count)
}

func politeWorker() {
	defer wg.Done()

	var count int

	for begin := time.Now(); time.Since(begin) <= runtime; {
		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()

		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()

		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()

		count += 1
	}

	fmt.Printf("Polite Boy died %d times\n", count)
}
