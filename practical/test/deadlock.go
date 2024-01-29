package main

import (
	"fmt"
	"sync"
	"time"
)

type x struct {
	mu sync.Mutex
	val int
}

// var wg sync.WaitGroup
func callmebaby (a *x, b *x) {
	defer wg.Done()

	a.mu.Lock()
	defer a.mu.Unlock()

	fmt.Printf("slept on %d\n", a.val)

	time.Sleep(2 * time.Second)

	fmt.Printf("woke up on %d\n", a.val)

	b.mu.Lock()
	defer b.mu.Unlock()

	fmt.Printf("Val: %d", 1);
}
