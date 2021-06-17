package main

import (
	"fmt"
	"sync"
)

func main() {
	var numCalcsCreated int
	//Cada elemento del pool es la direccion de un slice de bytes que ocupan 4K
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem // <1>
		},
	}

	// Seed the pool with 4KB
	//Creamos un pool con cuatro elementos
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			//Tomamos un elemento del pool. Hacemos un assert para comprobar que sea un puntero a un slice de bytes
			mem := calcPool.Get().(*[]byte) // <2>
			//Retornamos el elemento al pool
			defer calcPool.Put(mem)
		}()
	}

	wg.Wait()
	//Veamos cuantos elementos se han llegado a crear
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}
