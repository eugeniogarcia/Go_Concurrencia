package main

import (
	"fmt"
)

func main() {

	//Fuente de datos para nuestro pipeline
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		//El generador publica los datos por medio de un canal
		intStream := make(chan int)
		//La gorutina generará los datos
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					//termina la generación de datos
					return
					//Lógica de generación
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	//Pipeline que multiplica
	multiply := func(
		done <-chan interface{},
		intStream <-chan int,
		multiplier int,
	) <-chan int {
		//Define el canal para intercambiar el estado con la gorutina
		multipliedStream := make(chan int)
		//Define la gorutina
		go func() {
			//Señaliza el final del procesamiento cerrando el canal
			defer close(multipliedStream)
			//El procesamiento en un for-select
			for i := range intStream {
				select {
				case <-done:
					//Aborta el procesamiento
					return
					//Devuelbe el resultado del pipeline
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		//Retorna el canal
		return multipliedStream
	}

	//Pipeline que suma
	add := func(
		done <-chan interface{},
		intStream <-chan int,
		additive int,
	) <-chan int {
		//Define el canal para intercambiar el estado con la gorutina
		addedStream := make(chan int)
		//Define la gorutina
		go func() {
			//Señaliza el final del procesamiento cerrando el canal
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					//Aborta el procesamiento
					return
					//Devuelbe el resultado del pipeline
				case addedStream <- i + additive:
				}
			}
		}()
		//Retorna el canal
		return addedStream
	}

	//Canal para controlar la parada del pipeline - sino queremos esperar a que el generador deje de producir datos
	done := make(chan interface{})
	defer close(done)

	//Generador
	intStream := generator(done, 1, 2, 3, 4)
	//Lanza el pipeline. La respuesta del pipeline es un canal
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	//Procesamos los datos que vayan saliendo del canal
	for v := range pipeline {
		fmt.Println(v)
	}
}
