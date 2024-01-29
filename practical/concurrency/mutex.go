package main

import (
	"fmt"
	"sync"
)

func mutex() {
	count := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	incr := func ()  {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		count += 1
		fmt.Printf("Incrementing: %d\n", count)
	}

	decr := func ()  {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		count -= 1
		fmt.Printf("Decrementing: %d\n", count)
	}

	wg.Add(10)
	for i := 0; i < 5; i += 1 {
		go incr()
		go decr()
	}

	wg.Wait()
	fmt.Println("Arithmetic complete.")
}