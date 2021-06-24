package main

import (
	"fmt"
	"time"
)

func main() {
	//en done tenemos un canal de solo lectura que servira para señalizar cuando la go rutina tenga que terminar
	//strings es otro canal por el que pasamos los datos que queremos que procese la gorutina
	//Devolvemos un canal de lectura para señalizar cuando la gorutina ha terminado
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} { // <1>
		terminated := make(chan interface{})
		//lanza la go rutina que queremos controlar
		go func() {
			defer fmt.Println("doWork exited.")
			//Señalizamos que la gorutina termino
			defer close(terminated)
			for {
				select {
				//La gorutina recibe los datos que necesitamos procesar
				case s := <-strings:
					// Do something interesting
					fmt.Println(s)
				//Controlamos el fin de la ejecución de la gorutina con el canal done
				case <-done: // <2>
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() { // <3>
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated // <4>
	fmt.Println("Done.")
}
