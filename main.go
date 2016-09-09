package main

import (
	"fmt"
	"sync"
	"time"
)

func complexStuff(s string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("zi result is : %s", s)
}

func main() {
	MAX := 5
	jps := NewJobPubSub()
	var w sync.WaitGroup
	w.Add(MAX)
	//ps := pubsub.New(1)
	started := make(chan int, MAX)
	for i := 1; i <= MAX; i++ {
		go func(i int) {
			started <- i
			fmt.Printf("Starting routine %d \n", i)
			//arg := fmt.Sprintf("bouzin%d", i)
			arg := fmt.Sprintf("bouzin")
			val := jps.Compute(fmt.Sprintf(arg), func() interface{} { return complexStuff(arg) })
			fmt.Printf("goroutine %d says '%s' \n", i, val)
			w.Done()
		}(i)
	}
	w.Wait()

	arg := fmt.Sprintf("bouzin")
	val := jps.Compute(fmt.Sprintf(arg), func() interface{} { return complexStuff(arg) })
	fmt.Printf("main goroutine says '%s' \n", val)
}
