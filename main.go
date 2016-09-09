package main

import (
	"fmt"
	"github.com/cskr/pubsub"
	"sync"
	"time"
)

func main() {
	MAX := 5
	var w sync.WaitGroup
	w.Add(MAX)
	ps := pubsub.New(1)
	started := make(chan int, MAX)
	for i := 1; i <= MAX; i++ {
		go func(i int) {
			started <- i
			fmt.Printf("Starting routine %d \n", i)
			ch := ps.Sub("test")
			val := <-ch
			fmt.Printf("goroutine %d says %s \n", i, val)
			w.Done()
		}(i)
	}
	for len(started) != 5 {
		time.Sleep(10 * time.Millisecond)
	}
	ps.Pub("yes hello", "test")
	w.Wait()
}
