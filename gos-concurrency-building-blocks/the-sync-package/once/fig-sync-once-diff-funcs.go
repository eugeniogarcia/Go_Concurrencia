package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }

	var once sync.Once
	//Se ejecuta
	once.Do(increment)
	//Ya no se ejecuta porque la acabamos de ejecutar
	once.Do(decrement)

	fmt.Printf("Count: %d\n", count)
}
