# Job Pub/Sub

A mini library to leverage a pub/sub architecture to execute a job once and inform multiple go routines about the result. 
This prevents the same long-running job from being executed multiple times in parrallel as the multiple threads will be informed of the result once its computed.

## Usage

```
func complexStuff(s string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("The computed result is : %s", s)
}

func main() {
	jps := jobpubsub.NewJobPubSub()

	val := jps.Compute("complex-duty", func() interface{} {
		return complexStuff("hello")
	})
}
```

In the example above, you could call the compute method as many times as you want with the same ID (first parameter) during the second it computes and it will subscribe to a result channel and not execute the complexStuff method. See `job_pub_sub_test.go` for a full example.
