package main

import (
	"fmt"
	"time"
)

func chan_select() {
	start := time.Now()

	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()

	fmt.Println("Blocking on read...")

	select {
	case <-c:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

func chan_select_multiple() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 0; i < 1000; i += 1 {
		select {
		case <-c1:
			c1Count += 1
		case <-c2:
			c2Count += 1
		}
	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func chan_select_timeout() {
	var c <-chan int

	select {
	case <-c:
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}
}

func chan_select_default() {
	start := time.Now()
	var c1, c2 <-chan int

	select {

	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}
}

func chan_select_forloop() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0

loop:
	for {
		select {
		case <-done:
			break loop
		default: // without this, select will block till done streams or closes
		}

		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf(
		"Achieved %v cycles of work before signalled to stop.\n",
		workCounter,
	)
}
