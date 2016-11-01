package main

import (
	"fmt"
	"github.com/cskr/pubsub"
	"sync"
	"time"
)

type JobPubSub struct {
	lock      *sync.Mutex
	pubSubMap map[string]*pubsub.PubSub
}

func NewJobPubSub() *JobPubSub {
	jps := &JobPubSub{}
	jps.lock = &sync.Mutex{}
	jps.pubSubMap = make(map[string]*pubsub.PubSub)
	return jps
}

func (self *JobPubSub) Compute(id string, job func() interface{}) interface{} {

	c := self.getSubChan(id)
	// There is already someone taking care of it
	if c != nil {
		select {
		case val := <-c:
			return val
		case <-time.After(time.Second * 2):
			fmt.Println("Timeout reading from pubsub channel...")
			return job()
		}
	}

	// No one taking care of it so we compute it and publish it to any connected peer
	result := job()
	self.lock.Lock()
	defer self.lock.Unlock()
	self.pubSubMap[id].Pub(result, "result")
	self.pubSubMap[id] = nil
	return result
}

func (self *JobPubSub) getSubChan(id string) chan interface{} {
	self.lock.Lock()
	defer self.lock.Unlock()

	var c chan interface{}
	if self.pubSubMap[id] != nil {
		c = self.pubSubMap[id].Sub("result")
		fmt.Printf("Subscribing to existing computing for %s \n", id)
	} else {
		fmt.Printf("Creating new pubsub for %s \n", id)
		ps := pubsub.New(1)
		self.pubSubMap[id] = ps
	}
	return c
}
