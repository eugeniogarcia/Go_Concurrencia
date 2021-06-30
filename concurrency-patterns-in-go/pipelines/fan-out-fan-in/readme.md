En ocasiones podemos encontrarnos con que alguna de las stages del pipeline es especialmente costosa. En estos casos sería ideal poder paralelizar el procesamiento. Esto es lo que denominamos _fan-out_ y _fan-in_. Para poder hacer paralelizar una stage se tienen que dar las siguientes condiciones:

- Que sea una stage costosa
- Que el procesamiento de un dato no dependa de los datos que se han procesado antes, es decir, que haya independencia del orden en el que se procesan las cosas

Para demostrarlo creamos un pipeline que genera números primos. Los stages del pipeline son:

```go
repeatFn := func(done <-chan interface{},fn func() interface{},	) <-chan interface{}

take := func(done <-chan interface{},valueStream <-chan interface{},num int,) <-chan interface{}

toInt := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan int

primeFinder := func(done <-chan interface{}, intStream <-chan int) <-chan interface{}
```

Las primeras stages del pipeline generan enteros aleatorios:

```go
//Función que vamos a utilizar con el generador
rand := func() interface{} { return rand.Intn(50000000) }
//Generador
randIntStream := toInt(done, repeatFn(done, rand))
```

La siguiente stage la procesaremos en paralelo - creamos `runtime.NumCPU()` stages en paralelo:

```go
	//Crea el fan-out
	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	//Creamos un slice de canales con un tamaño numFinders 
    finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Primes:")
	for i := 0; i < numFinders; i++ {
        //Cada elemento del slice es una stage en paralelo. Cada una procesa uno de los enteros
		finders[i] = primeFinder(done, randIntStream)
	}
```

Con el _fan-in_ consolidamos cada una de las stages que se ejecutan en paralelo

```go
	//Crea el fan-in
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}
```

## Fan-in

Con la stage fan-in lo que hacemos es tomar la salida de varias stages de entrada - que se están procesando en paralelo -, y la coloca en el canal de salida.

- Tomamos un variant de canales

```go
channels ...<-chan interface{}
```

- y devuelbe un solo canal

```go
) <-chan interface{}
```

- La stage termina cuando hayan terminado todas las stages de entrada. Para coordinar cada entrada usamos un `WaitGroup`

```go
var wg sync.WaitGroup // <2>

wg.Add(len(channels)) // <4>
```

- El procesamiento de la stage se hace tambien en una gorutina. Como tenemos varias entradas a procesar, se lanzarán varias gorutinas:

```go
multiplex := func(c <-chan interface{}) { // <3>
    //Notificamos que se ha terminado la stage
    defer wg.Done()
    //Tomamos los datos de cada una de las n-stages de entrada y los volcamos al único canal de salida 
    for i := range c {
        select {
        case <-done:
            return
        case multiplexedStream <- i:
        }
    }
}
```

```go
for _, c := range channels {
    go multiplex(c)
}
```

- Coordinar que cuando todas las entradas hayan terminado se termine también la stage se hace en otra gorutina:

```go
go func() { // <5>
    wg.Wait()
    close(multiplexedStream)
}()
```

Combinando todas estas partes, el _fan-in_ queda como sigue:

```go
//Implementa el fan-in. Toma un variant de canales de entrada
fanIn := func(
    done <-chan interface{},
    channels ...<-chan interface{},
) <-chan interface{} { // <1>
    var wg sync.WaitGroup // <2>
    //Como cualquier otra stage se comunica con el exterior con un canal
    multiplexedStream := make(chan interface{})

    //Esta es la lógica que se aplicará a cada uno de los canales a procesar en paralelo
    multiplex := func(c <-chan interface{}) { // <3>
        //Notificamos que se ha terminado la stage
        defer wg.Done()
        //Tomamos los datos de cada una de las n-stages de entrada y los volcamos al único canal de salida 
        for i := range c {
            select {
            case <-done:
                return
            case multiplexedStream <- i:
            }
        }
    }

    // Indicamos cuantas stages paralelas tenemos que procesar antes poder terminar la stage fan-in
    wg.Add(len(channels)) // <4>
    //Procesa en una gorutina los resultados de cada stage de entrada
    for _, c := range channels {
        go multiplex(c)
    }

    // Cuando se han terminado todas las stages de entrada, terminamos la fan-in
    go func() { // <5>
        wg.Wait()
        close(multiplexedStream)
    }()

    return multiplexedStream
}
```