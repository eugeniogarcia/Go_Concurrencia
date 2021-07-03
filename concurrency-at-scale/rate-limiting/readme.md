Usamos el paquete `golang.org/x/time/rate`.

El tipo `rate.Limit` nos dice cuantas peticiones/eventos por segundo vamos a admitir. También podemos crear un `Limit` con este helper:

```go
rate.Every(duration / time.Duration(eventCount))
```

Este helper nos permite definir un Limit indicando como argumento el tiempo que tiene que pasar entre eventos.

Una vez tenemos un `Limit` ya podemos crear el _Rate Limiter_. Con este helper creamos el rate limitir, indicando un `Limit` y un _burstable_:

```go
secondLimit := rate.NewLimiter(rate.Every(8), 1)
```

Por último cuando queramos controlar el uso de un recurso lo que tendremos que hacer es usar `Wait` con el contexto. Wait bloqueará la ejecución sino hay un `token` disponible:

```go
func (a *APIConnection) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil { // <2>
		return err
	}
```

`Wait` es un shortcut que equivale a `WaitN(ctx,1)`.


