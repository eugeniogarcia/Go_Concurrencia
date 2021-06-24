Cuando tenemos una gorutina que procesa trabajo en un loop - por ejemplo con un _for select loop_ -, necesitamos un mecanismo para poder señalizar a la gorutina cuando tiene que terminar su trabajo. Para hacer esta comunicación usaremos un canal de solo lectura.

```go
done := make(chan interface{})
```

Pasaremos este canal como primer argumento - es una convención - a la gorutina que queramos controlar:

```go
terminated := doWork(done, nil)
```

```go
	//en done tenemos un canal de solo lectura que servira para señalizar cuando la go rutina tenga que terminar
	//strings es otro canal por el que pasamos los datos que queremos que procese la gorutina
	//Devolvemos un canal de lectura para señalizar cuando la gorutina ha terminado
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} { // <1>
		terminated := make(chan interface{})
		//lanza la go rutina que queremos controlar
		go func() {
			defer fmt.Println("doWork exited.")
			//Señalizamos que la gorutina termino
			defer close(terminated)
			for {
				select {
				//La gorutina recibe los datos que necesitamos procesar
				case s := <-strings:
					// Do something interesting
					fmt.Println(s)
				//Controlamos el fin de la ejecución de la gorutina con el canal done
				case <-done: // <2>
					return
				}
			}
		}()
		return terminated
	}
```

Además del canal _done_ de solo lectura que usamos para avisar a la gorutina de cuando tiene que terminar, tenemos otros dos canales en este ejemplo. Uno que usamos como input a la go rutina, de solo lectura, que proporciona el mecanismo con el cual alimentar de datos para que la gorutina trabaje. El otro canal, también de lectura, es devuelto por la gorutina con la finalidad de comunicar a quien este interesado de cuando se ha terminado el procesamiento.