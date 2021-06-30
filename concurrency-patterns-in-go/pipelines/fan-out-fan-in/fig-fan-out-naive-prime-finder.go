package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}
	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}
	toInt := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case intStream <- v.(int):
				}
			}
		}()
		return intStream
	}
	primeFinder := func(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
		primeStream := make(chan interface{})
		go func() {
			defer close(primeStream)
			for integer := range intStream {
				integer -= 1
				prime := true
				for divisor := integer - 1; divisor > 1; divisor-- {
					if integer%divisor == 0 {
						prime = false
						break
					}
				}

				if prime {
					select {
					case <-done:
						return
					case primeStream <- integer:
					}
				}
			}
		}()
		return primeStream
	}

	//Implementa el fan-in. Toma un variant de canales de entrada
	fanIn := func(
		done <-chan interface{},
		channels ...<-chan interface{},
	) <-chan interface{} { // <1>
		var wg sync.WaitGroup // <2>
		//Como cualquier otra stage se comunica con el exterior con un canal
		multiplexedStream := make(chan interface{})

		//Esta es la lógica que se aplicará a cada uno de los canales a procesar en paralelo
		multiplex := func(c <-chan interface{}) { // <3>
			//Notificamos que se ha terminado la stage
			defer wg.Done()
			//Tomamos los datos de cada una de las n-stages de entrada y los volcamos al único canal de salida
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		// Indicamos cuantas stages paralelas tenemos que procesar antes poder terminar la stage fan-in
		wg.Add(len(channels)) // <4>
		//Procesa en una gorutina los resultados de cada stage de entrada
		for _, c := range channels {
			go multiplex(c)
		}

		// Cuando se han terminado todas las stages de entrada, terminamos la fan-in
		go func() { // <5>
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	done := make(chan interface{})
	defer close(done)

	//Medimos cuanto tiempo se precisa
	start := time.Now()

	//Función que vamos a utilizar con el generador
	rand := func() interface{} { return rand.Intn(50000000) }
	//Generador
	randIntStream := toInt(done, repeatFn(done, rand))

	//Crea el fan-out
	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	//Creamos un slice de canales con un tamaño numFinders
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Primes:")
	for i := 0; i < numFinders; i++ {
		//Cada elemento del slice es una stage en paralelo. Cada una procesa uno de los enteros
		finders[i] = primeFinder(done, randIntStream)
	}

	//Crea el fan-in
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	//Paramos el cronometro
	fmt.Printf("Search took: %v", time.Since(start))
}
