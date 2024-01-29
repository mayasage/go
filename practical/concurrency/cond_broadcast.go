package main

import (
	"fmt"
	"sync"
)

func condBroadcast() {
	type Button struct {
		Clicked *sync.Cond
	}
	btn := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func (c *sync.Cond, fn func())  {
		var wg sync.WaitGroup
		wg.Add(1)

		go func ()  {
			// fmt.Printf("registered\n") // no deadlock - Why?
			wg.Done()
			// fmt.Printf("registered\n") // deadlock - Why?
			c.L.Lock()
			defer c.L.Unlock()
			// fmt.Printf("registered\n") // deadlock - Why?
			c.Wait()
			fn()
		}()

		wg.Wait()
	}

	var wg sync.WaitGroup

	wg.Add(3)
	subscribe(btn.Clicked, func ()  {
		fmt.Println("Maximizing window.")
		wg.Done()
	})
	subscribe(btn.Clicked, func ()  {
		fmt.Println("Displaying annoying dialog box!")
		wg.Done()
	})
	subscribe(btn.Clicked, func ()  {
		fmt.Println("Mouse clicked.")
		wg.Done()
	})

	btn.Clicked.Broadcast()
	wg.Wait()
}