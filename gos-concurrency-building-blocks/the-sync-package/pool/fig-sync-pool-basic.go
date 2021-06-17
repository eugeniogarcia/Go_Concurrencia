package main

import (
	"fmt"
	"sync"
)

func main() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance.")
			return struct{}{}
		},
	}

	//Como no hay nada en el pool, se creara una instancia
	myPool.Get() // <1>
	//Como no hay nada en el pool, se creara una instancia
	instance := myPool.Get() // <1>
	//devolvemos la instancia al pool
	myPool.Put(instance) // <2>
	//Recuperamos la instancia del pool
	myPool.Get() // <3>
}
