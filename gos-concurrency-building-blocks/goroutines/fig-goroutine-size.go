package main

import (
	"fmt"
	"runtime"
	"sync"
)

/*
Demuestra que las gorutinas son componentes muy ligeros, que ocupan poca memoria.

Vemos el espacio de memoria ocupado antes y después de haber lanzado 10000 gorutinas - que no hacen nada y quedan bloqueadas, sin terminar.
*/
func main() {
	memConsumed := func() uint64 {
		runtime.GC()
		//Obtenermos la memoria ocupada
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	//Esperamos a que nos llegue algo por el canal, lo que en terminos prácticos significa que se bloqueara la ejecución
	noop := func() { wg.Done(); <-c } // <1>

	const numGoroutines = 1e4 // <2>
	wg.Add(numGoroutines)
	before := memConsumed() // <3>
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed() // <4>
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}
