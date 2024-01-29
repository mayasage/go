package main

import (
	"fmt"
	"sync"
)

func syncAdd() {
	var wg sync.WaitGroup

	hello := func (id int) {
		defer wg.Done()
		fmt.Printf("Hello from %v!\n", id)
	}

	const numGreeters = 5
	wg.Add(numGreeters)
	for i := 0; i < numGreeters; i += 1 {
		go hello(i + 1)
	}
	wg.Wait()
}
