package main

import (
	"fmt"
	"sync"
)

func pool() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance.")
			return struct{}{}
		},
	}

	myPool.Get()
	inst := myPool.Get()
	myPool.Put(inst)
	myPool.Get()
}

func pool2() {
	var count int
	myPool := &sync.Pool{
		New: func() interface{} {
			count += 1
			byteArr := make([]byte, 1024)
			return &byteArr
		},
	}

	myPool.Put(myPool.Get())
	myPool.Put(myPool.Get())
	myPool.Put(myPool.Get())
	myPool.Put(myPool.Get())

	const numWorkers = 1024 * 1024

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i += 1 {
		go func ()  {
			defer wg.Done()
			mem := myPool.Get().(*[]byte)
			defer myPool.Put(mem)
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", count)
}
