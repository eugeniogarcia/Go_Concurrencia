package main

import (
	"sync"
)

func main() {
	var onceA, onceB sync.Once
	var initB func()
	initA := func() {
		println("A")
		onceB.Do(initB)
	}
	initB = func() {
		println("B")
		onceA.Do(initA)
	} // <1>
	println("prepara")
	onceA.Do(initA)
	println("termina") // <2>
}
