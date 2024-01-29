package main

import (
	"fmt"
	"sync"
)

func once() {
	var count int

	incr := func() {
		count += 1
	}

	var once sync.Once

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i += 1 {
		go func ()  {
			defer wg.Done()
			once.Do(incr)
		}()
	}
	wg.Wait()
	fmt.Printf("Count is %d\n", count)
}

func once2() {
	var count int

	incr := func ()  {
		count += 1
	}
	decr := func ()  {
		count -= 1
	}

	var once sync.Once
	once.Do(incr)
	once.Do(decr)

	fmt.Printf("Count: %d\n", count)
}

func onceDeadlock() {
	var onceA sync.Once
	var initA, initB func()

	initA = func ()  {
		onceA.Do(initB)
	}
	initB = func ()  {
		onceA.Do(initA)
	}

	onceA.Do(initA)
}