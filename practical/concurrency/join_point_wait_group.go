package main

import (
	"fmt"
	"sync"
)

func joinPointWaitGroup () {
	var wg sync.WaitGroup

	sayHello := func () {
		defer wg.Done()
		fmt.Printf("hello")
	}

	wg.Add(1)
	go sayHello()
	wg.Wait()
}
