package main

import (
	"fmt"
	"sync"
)

func closure() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func ()  {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()
}

func closure2() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func (salutation string)  {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()
}