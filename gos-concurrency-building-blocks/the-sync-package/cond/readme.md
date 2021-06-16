The comment for the Cond type really does a great job of describing its purpose:

...a rendezvous point for goroutines waiting for or announcing the occurrence of an event.

In that definition, an “event” is any arbitrary signal between two or more goroutines that carries no information other than the fact that it has occurred. Very often you’ll want to wait for one of these signals before continuing execution on a goroutine.

Crea la condición:

```go
c := sync.NewCond(&sync.Mutex{})
```

Entramos en la sección crítica:

```go
c.L.Lock()
```

Salimos de la sección crítica:

```go
c.L.Unlock()
```

Hasta aquí todo "normal", tenemos las mismas capacidades que tendríamos con un _Mutex_. Ahora viene la parte interesante, la que __permite emitir y recibir señales entre go rutinas__. 

Con _Wait_ indicamos que necesitamos ser notificados - __recibir__ . __Cuando llamamos a wait, por detras se hace un Unlock, y la go rutina se suspende__. El runtime scheduler mantiene una relación de go rutinas que se han suspendido - y que están esperando a ser despertadas:

```go
c.Wait()
```

Para emitir tenemos dos opciones:
- Signal. Signal emitirá una señal al runtime para indicar que debe notificarse a otra go rutina que estuviera suspendida, para que retome su trabajo. Cuando usamos Signal lo que hara el runtime es despertar la go rutina que lleve más tiempo esperando. 
-Broadcast. Con Broadcast, se notificará a todas a la vez, y se despertara una de ellas.

```go
c.Signal()
```

```go
c.Broadcast()
```