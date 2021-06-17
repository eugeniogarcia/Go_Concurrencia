package main

import (
	"fmt"
)

func main() {
	intStream := make(chan int)
	close(intStream)
	//Cuando leemos de un canal, en realidad podemos recibir un par de datos, el valor y un booleano. El booleano nos indica si el canal esta o no abierto. En este caso recibiremos el valor por defecto, que en un integer es 0, y false porque el canal esta cerrado
	integer, ok := <-intStream // <1>
	fmt.Printf("(%v): %v", ok, integer)
}
