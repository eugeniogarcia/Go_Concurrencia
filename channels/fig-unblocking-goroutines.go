package main

import (
	"fmt"
	"sync"
)

func main() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin // <1>
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	//Al cerrar el canal, todas aquellas go rutinas que estuvieran bloqueadas esperando datos poe el canal, se desbloquearan - y recibiran el valor por defecto
	close(begin) // <2>
	wg.Wait()
}
