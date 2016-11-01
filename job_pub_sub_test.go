package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func complexStuff(s string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("zi result is : %s", s)
}

func TestJobPubSub(t *testing.T) {
	count := 0
	m := &sync.Mutex{}
	MAX := 5
	jps := NewJobPubSub()
	var w sync.WaitGroup
	w.Add(MAX)
	for i := 1; i <= MAX; i++ {
		go func(i int) {
			//fmt.Printf("Starting routine %d \n", i)
			arg := "bouzin"
			val := jps.Compute("complex-"+arg, func() interface{} {
				m.Lock()
				count = 1
				//fmt.Printf("Count is %d \n", count)
				m.Unlock()
				return complexStuff(arg)
			})
			if val != "zi result is : bouzin" {
				t.Error(fmt.Sprintf("Goroutine %d didn't return the proper value (%s)", i, val))
			}
			//fmt.Printf("goroutine %d says '%s' \n", i, val)
			w.Done()
		}(i)
	}
	w.Wait()
	if count != 1 {
		t.Error(fmt.Sprintf("Inner method was executed %d times instead of 1", count))
	}

}
