// A mini library to leverage a pub/sub architecture to execute a job once and inform multiple go routines about the result.
// This prevents the same long-running job from being executed multiple times in parrallel as the multiple threads will be informed of the result once its computed.
package jobpubsub

import (
	"fmt"
	"github.com/cskr/pubsub"
	"sync"
	"time"
)

type JobPubSub struct {
	lock      *sync.Mutex
	pubSubMap map[string]*pubsub.PubSub
	Timeout   int
}

// NewJobPubSub creates Pub/Sub struct which can then be used to compute data from
func NewJobPubSub() *JobPubSub {
	jps := &JobPubSub{}
	jps.lock = &sync.Mutex{}
	jps.pubSubMap = make(map[string]*pubsub.PubSub)
	jps.Timeout = 1
	return jps
}

// Compute will execute a function that takes an interface as the only argument and returns an interface
// It must be given an identifier of the job (the function) that is being executed
// In the event a job is running with the same ID, it will subscribe to the return of that other job
// The return value will be the same for any jobs running concurrently with the same ID
// JobPubSub.Timeout can be used to specify a timeout (1 second by default, -1 to deactivate timeout)
func (self *JobPubSub) Compute(id string, job func() interface{}) interface{} {

	c := self.getSubChan(id)
	// There is already someone taking care of it
	if c != nil {
		select {
		case val := <-c:
			return val
		case <-time.After(time.Second * self.Timeout):
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

// getSubChan will return a subscriber channel for a given ID
func (self *JobPubSub) getSubChan(id string) chan interface{} {
	self.lock.Lock()
	defer self.lock.Unlock()

	var c chan interface{}
	if self.pubSubMap[id] != nil {
		c = self.pubSubMap[id].Sub("result")
		//fmt.Printf("Subscribing to existing computing for %s \n", id)
	} else {
		//fmt.Printf("Creating new pubsub for %s \n", id)
		ps := pubsub.New(1)
		self.pubSubMap[id] = ps
	}
	return c
}
