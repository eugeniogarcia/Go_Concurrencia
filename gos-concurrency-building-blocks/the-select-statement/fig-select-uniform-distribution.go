package main

import (
	"fmt"
)

//Con este ejemplo se demuestra que el select no da preferencia a un canal sobre otro. Cuando hay datos disponibles en más de un canal, select elige de forma aleatoria de que canal leer
func main() {
	c1 := make(chan interface{})
	//Cerramos el canal. De esta forma siempre habrá datos disponibles cuando hagamos el select
	close(c1)
	c2 := make(chan interface{})
	//Cerramos el canal. De esta forma siempre habrá datos disponibles cuando hagamos el select
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}
