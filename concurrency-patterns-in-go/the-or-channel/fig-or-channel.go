package main

import (
	"fmt"
	"time"
)

//Si tenemos varios canales que señalan que el trabajo se ha completado, con esta función los agrupamos en un canal que actua como OR. El canal nos señalizará el fin cuando alguno de los canales se haya completado

func main() {
	//Define una variable de tipo función, que adminte un variant de canales de salida y que retorna un canal de salida
	var or func(channels ...<-chan interface{}) <-chan interface{}
	//Instancia la variable or con una función
	or = func(channels ...<-chan interface{}) <-chan interface{} { // <1>
		switch len(channels) {
		case 0: // <2>
			return nil
		case 1: // <3>
			return channels[0]
		}
		//Si hay más de dos canales en la entrada, este canal que definimos aquí los "une" a todos
		orDone := make(chan interface{})
		go func() { // <4>
			//Se ha terminado alguno de los n-canales
			defer close(orDone)

			switch len(channels) {
			case 2: // <5>
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default: // <6>
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				//La parte recursiva de la solución
				case <-or(append(channels[3:], orDone)...): // <6>
				}
			}
		}()
		return orDone
	}
	sig := func(after time.Duration) <-chan interface{} { // <1>
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now() // <2>
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start)) // <3>
}
