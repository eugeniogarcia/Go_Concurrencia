package main

import (
	"fmt"
	"time"
)

// Demuestra una técnica para definir un time-out usando un select. Usamos time.After que publica en un canal un pulso después del tiempo que le indiquemos, en este ejemplo, tras un segundo
func main() {
	var c <-chan int
	select {
	case <-c: // <1>
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}
}
