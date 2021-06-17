Crea un pool de recursos. Lo aconsejable es que los recursos sean todos ellos homogeneos.

El pool se inicializa especificando cual es el constructor que se tiene que usar para crear un recurso:

```go
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem // <1>
		},
	}
```

A partir de este momento podemos ya utilizar el pool. Para obtener un recurso:

```go
mem := calcPool.Get()
```

Este método recupera un recurso del pool, o si no hay recursos disponibles, crea un recurso - usando el constructor. 

Cuando ya no necesitamos usar más el recurso, lo devolvemos al pool:

```go
defer calcPool.Put(mem)
```

Podemos inicializar el pool creando algunos recursos con Put, antes de que alguien los necesite:

```go
	//Creamos un pool con cuatro elementos
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
```
