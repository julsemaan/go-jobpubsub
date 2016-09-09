package main

import (
	"fmt"
	"github.com/cskr/pubsub"
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var pubSubMap = make(map[string]*pubsub.PubSub)

func compute(id string, job func() interface{}) interface{} {

	c := getSubChan(id)
	// There is already someone taking care of it
	if c != nil {
		select {
		case val := <-c:
			return val.(string)
		case <-time.After(time.Second * 2):
			fmt.Println("Timeout reading from pubsub channel...")
			return job()
		}
	}

	// No one taking care of it so we compute it and publish it to any connected peer
	result := job()
	lock.Lock()
	pubSubMap[id].Pub(result, "result")
	pubSubMap[id] = nil
	lock.Unlock()
	return result
}

func getSubChan(id string) chan interface{} {
	lock.Lock()
	defer lock.Unlock()

	var c chan interface{}
	if pubSubMap[id] != nil {
		c = pubSubMap[id].Sub("result")
		fmt.Printf("Subscribing to existing computing for %s \n", id)
	} else {
		fmt.Printf("Creating new pubsub for %s \n", id)
		ps := pubsub.New(1)
		pubSubMap[id] = ps
	}
	return c
}

func complexStuff(s string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("zi result is : %s", s)
}

func main() {
	MAX := 5
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
			val := compute(fmt.Sprintf(arg), func() interface{} { return complexStuff(arg) })
			fmt.Printf("goroutine %d says '%s' \n", i, val)
			w.Done()
		}(i)
	}
	w.Wait()
}
