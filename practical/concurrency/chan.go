package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

func chan_test() {
	stringStream := make(chan string)

	go func() {
		stringStream <- "Hello channels!"
	}()

	salutation, ok := <-stringStream
	fmt.Printf("(%v): %v", ok, salutation)
}

func chan_readOnly() {
	// readStream := make(<-chan interface{})
	// readStream <- struct{}{} // This write to a read-only stream
}

func chan_writeOnly() {
	// writeStream := make(chan<- interface{})
	// <- writeStream // Can't read from a write-only stream
}

func chan_closed() {
	intStream := make(chan int)
	close(intStream)
	integer, ok := <- intStream
	fmt.Printf("(%v): %v", ok, integer)
}

func chan_itr() {
	intStream := make(chan int)

	go func() {
		defer close(intStream)

		for i := 1; i <= 5; i+= 1 {
			intStream <- i
		}
	}()

	for intVal := range intStream {
		fmt.Printf("%v ", intVal)
	}
}

func chan_unblockAll() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i += 1 {
		go func (i int)  {
			defer wg.Done()
			<- begin
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()
}

func chan_buffered() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i += 1 {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for intVal := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %v.\n", intVal)
	}
}

func chan_nil() {
	// panic - reading from a nil chan
	// var dataStream chan interface{}
	// <-dataStream

	// panic - writing to a nil chan
	// var dataStream chan interface{}
	// dataStream <- struct{}{}

	// panic - closing a nil chan
	// var dataStream chan interface{}
	// close(dataStream)
}

func chan_ownership() {
	chanOwner := func () <-chan int  {
		resStream := make(chan int, 5)

		go func ()  {
			defer close(resStream)
			for i := 0; i < 5; i += 1 {
				resStream <- i
			}
		}()

		return resStream
	}

	resStream := chanOwner()
	for intVal := range resStream {
		fmt.Printf("Received: %d\n", intVal)
	}
	fmt.Println("Done receiving!")
}