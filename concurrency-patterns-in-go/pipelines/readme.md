# Pipelines (fig-pipelines-chan-stream-processing)

Demuestra como definir un pipeline. Un pipeline es un conjunto de stages; Una stage es una función que realiza un trabajo concreto, que toma unos datos de entrada y genera unos de salida. Los datos de entrada son inmutables.

Vamos a construir los stages usando canales para intercambiar datos de una a otra función.

Un stage sigue el siguiente patrón:

```go
	//Pipeline que multiplica
	multiply := func(
		done <-chan interface{},
		intStream <-chan int,
		multiplier int,
	) <-chan int {
        //Define el canal para intercambiar el estado con la gorutina
		multipliedStream := make(chan int)
        //Define la gorutina
		go func() {
            //Señaliza el final del procesamiento cerrando el canal
			defer close(multipliedStream)
            //El procesamiento en un for-select
			for i := range intStream {
				select {
				case <-done:
                    //Aborta el procesamiento
					return
                //Devuelbe el resultado del pipeline
				case multipliedStream <- i * multiplier:
				}
			}
		}()
        //Retorna el canal
		return multipliedStream
	}
```

Son una función que toma como entrada un canal y opcionalmente otra serie de datos, y devuelve un canal:

```go
	multiply := func(
		done <-chan interface{},
		intStream <-chan int,
		multiplier int,
	) <-chan int {
```

Notese que los canales son de solo lectura y que están tipificados - en nuestro caso, intercambian `int`. También usamos un canal `done` para controlar cuando deseamos detener el procesamiento. El procesamiento se hace en una gorutina, de modo que el procesamiento de un pipeline se hace de forma concurrente:

```go
//Define el canal para intercambiar el estado con la gorutina
multipliedStream := make(chan int)

//Define la gorutina
go func() {
    //Señaliza el final del procesamiento cerrando el canal
    defer close(multipliedStream)
    //El procesamiento en un for-select
    for i := range intStream {
        select {
        case <-done:
            //Aborta el procesamiento
            return
        //Devuelbe el resultado del pipeline
        case multipliedStream <- i * multiplier:
        }
    }
}()

//Retorna el canal
return multipliedStream
```

## Generador

Para utilizar el pipeline necesitamos una fuente de datos, un canal, que podamos usar para limentar el pipeline: un generador:

```go
	//Fuente de datos para nuestro pipeline
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		//El generador publica los datos por medio de un canal
		intStream := make(chan int)
		//La gorutina generará los datos
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					//termina la generación de datos
					return
					//Lógica de generación
				case intStream <- i:
				}
			}
		}()
		return intStream
	}
```

El generador sigue también un patrón parecido al del pipeline, con la salvedad de que no tiene como entrada un canal. Por lo demás es bastante auto-explicativo

## Uso

Vamos a ver como usar un generador para crear la fuente de datos, y luego procesarlos en un pipeline:

```go
	//Canal para controlar la parada del pipeline - sino queremos esperar a que el generador deje de producir datos
	done := make(chan interface{})
	defer close(done)

	//Generador
	intStream := generator(done, 1, 2, 3, 4)
	//Lanza el pipeline. La respuesta del pipeline es un canal
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	//Procesamos los datos que vayan saliendo del canal
	for v := range pipeline {
		fmt.Println(v)
	}
```

- El generador devuelve un canal por el que se irán publicando los datos que el pipeline tiene que procesar
- El pipeline combina varios stages
- La salida de cada stage, y del pipeline como tal, es un canal. Procesamos el canal con `range`
- Las diferentes stages del pipeline se irán procesando de forma concurrente
- Podemos detener la ejecución del pipeline - y del generador - cerrando el canal `done`

# Ejemplos (fig-utilizing-string-stage)

## Generador Repeat

Vamos a ver como podríamos crear un generador que repita la emisión de un juego de datos:

- Pasaremos como input un variant de un tipo genérico

```go
values ...interface{},
```

- Obtenermos los datos del variant - el variant no deja de ser un slice -, y cuando les hemos extraido todos, volvemos a empezar:

```go
for {
    for _, v := range values {
        select {
        case <-done:
            return
        case valueStream <- v:
        }
    }
}
```

Todo combinado:

```go
	repeat := func(
		done <-chan interface{},
		values ...interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}
```

## Repeat de una función

Similar al caso anterior, este generador genera de forma repetida un juego de datos. La diferencia con el ejemplo anterior es que los datos son generados con la aplicación de una función. El parametro del generador es una función:

```go
fn func() interface{},
```

Los datos que se generan son los que produce la función:

```go
case valueStream <- fn():
```

El generador quedaría como sigue:

```go
	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}
```

Por ejemplo, un generador de números aleatorios podría ser el que sigue:

```go
rand := func() interface{} { return rand.Intn(50000000) }

randIntStream := toInt(done, repeatFn(done, rand))
```

## Take

Veamos una stage que toma un número determinado de datos del canal y luego termina.

- Los datos de entrada del canal vuelven a ser genéricos:

```go
valueStream <-chan interface{}
```

El stage tiene como argumento el número de valores a tomar del canal:

```go
num int,
```

Tomamos el número de valores del canal, y terminamos el stage:

```go
go func() {
    defer close(takeStream)
    for i := 0; i < num; i++ {
        select {
        case <-done:
            return
        case takeStream <- <-valueStream:
        }
    }
}()
```

Todo combinado:

```go
	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}
```

## toString

Con `take` y `repeat` hemos usado un tipo genérico - `interface{}`. Supongamos que en un momento determinado del pipeline quisieramos convertir el datos a string. Podríamos usar un stage como este:

- El stage toma un tipo genérico, pero devuelve un string:

```go
toString := func(done <-chan interface{},valueStream <-chan interface{},) <-chan string {
```

- La gortina del stage se limita a hacer la assertion:

```go
case stringStream <- v.(string):
```

Todo combinado:

```go
	toString := func(
		done <-chan interface{},
		valueStream <-chan interface{},
	) <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case stringStream <- v.(string):
				}
			}
		}()
		return stringStream
	}
```
