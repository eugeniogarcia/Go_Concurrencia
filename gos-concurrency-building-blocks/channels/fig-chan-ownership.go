package main

import (
	"fmt"
)

/*
The first thing we should do to put channels in the right context is to assign channel
ownership. I’ll define ownership as being a goroutine that instantiates, writes, and
closes a channel. Much like memory in languages without garbage collection, it’s
important to clarify which goroutine owns a channel in order to reason about our
programs logically.

Unidirectional channel declarations are the tool that will allow us
to distinguish between goroutines that own channels and those that only utilize them:
channel owners have a write-access view into the channel (chan or chan<-), and
channel utilizers only have a read-only view into the channel (<-chan).
*/
func main() {
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5) // <1>
		go func() {                       // <2>
			defer close(resultStream) // <3>
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream // <4>
	}

	resultStream := chanOwner()
	for result := range resultStream { // <5>
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
}
