package main

import (
	"fmt"
	"math/rand"
)

func main() {
	//Simula la gorutina que hace trabajo
	//Devolbemos un canal con el heartbeat, y otro con los datos, en este caso int
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		//Los canales para el heartbeat y los datos
		heartbeatStream := make(chan interface{}, 1) // <1>
		workStream := make(chan int)

		go func() {
			defer close(heartbeatStream)
			defer close(workStream)

			//Envia un heartbeat
			for i := 0; i < 10; i++ {
				select { // <2>
				case heartbeatStream <- struct{}{}:
				default: // <3>
				}

				//Realiza el trabajo
				select {
				case <-done:
					return
				case workStream <- rand.Intn(10):
				}
			}
		}()

		return heartbeatStream, workStream
	}

	done := make(chan interface{})
	defer close(done)

	heartbeat, results := doWork(done)

	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}
