package main

import (
	"fmt"
)

func main() {
	intStream := make(chan int)
	go func() {
		defer close(intStream) // <1>
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()

	//Usando range estaremos recibiendo valores hasta que el canal se cierre. En el momento que se cierre saldremos de este loop
	for integer := range intStream { // <2>
		fmt.Printf("%v ", integer)
	}
}
