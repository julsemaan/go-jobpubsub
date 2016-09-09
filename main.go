package main

import (
	"fmt"
	"github.com/cskr/pubsub"
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var pubSubMap = make(map[string]*pubsub.PubSub)

func compute(s string) string {

	c := getSubChan(s)
	if c != nil {
		val := <-c
		return val.(string)
	}

	result := complexStuff(s)
	lock.Lock()
	pubSubMap[s].Pub(result, "result")
	pubSubMap[s] = nil
	lock.Unlock()
	return result
}

func getSubChan(s string) chan interface{} {
	lock.Lock()
	defer lock.Unlock()

	var c chan interface{}
	if pubSubMap[s] != nil {
		c = pubSubMap[s].Sub("result")
		fmt.Printf("Subscribing to existing computing for %s \n", s)
	} else {
		fmt.Printf("Creating new pubsub for %s \n", s)
		ps := pubsub.New(1)
		pubSubMap[s] = ps
	}
	return c
}

func complexStuff(s string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("zi result is : %s", s)
}

func main() {
	MAX := 20
	var w sync.WaitGroup
	w.Add(MAX)
	//ps := pubsub.New(1)
	started := make(chan int, MAX)
	for i := 1; i <= MAX; i++ {
		go func(i int) {
			started <- i
			fmt.Printf("Starting routine %d \n", i)
			val := compute(fmt.Sprintf("bouzin"))
			fmt.Printf("goroutine %d says '%s' \n", i, val)
			w.Done()
		}(i)
	}
	w.Wait()
}
