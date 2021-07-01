Con el paquete `context` creamos un contexto. El contexto nos permite dos cosas:

- Proporcionar un mecanismo para cancelar stages en el pipeline (el padre puede cancelar la ejecución)
- Proporcionar un data-bag para transmitir datos comunes a lo largo de todas las stafes del pipeline

Para crear un contexto:

```go
context.Background()
```

El contexto es inmutable. No podemos cambiar el data-bag ni cambiar forzar la conclusión de la ejecución. Si queremos, por ejemplo, controlar la ejecución aguas abajo, decidir cuando abortarla, tendremos que usar una de las funciones de _context_, y generar un nuevo contexto con la función que nos permitirá controlarlo. Por ejemplo, si una stage quiere poder controlar cuando se terminará el procesamiento de las stages hijas, hará esto con el contexto que haya recibido:

```go
ctx, cancel := context.WithCancel(context.Background())
```

El contexto que se pasará a los hijos será `ctx`, y cuando quiera cancelar la ejecución, llamará a:

```go
cancel()
```

Lo mismo sucede con el data-bag. Como context es inmutable se a de seguir el siguiente patron:

```go
ctx := context.WithValue(context.Background(), "userID", userID)
ctx = context.WithValue(ctx, "authToken", authToken)
```

Podemos también añadir un time-out a un contexto:

```go
ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
```

Cuando un contexto tiene un time-out, el consumidor del contexto podrá interrogarle para saber si hay un time-out y cual es el valor:

```go
if deadline, ok := ctx.Deadline(); ok {
    if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
```

Aquí cuando `ok` es `true`significa que hay un deadline fijado, y en `deadline` tenemos cual es valor.

De la misma forma que podemos interrogar a un contexto para saber si tiene un deadline, podemos recuperar datos del data-bag:

```go
ctx.Value("userID"),
ctx.Value("authToken"),
```

Podemos comprobar si el contexto a finalizado, recuperando el canal `done`:

```go
func locale(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err() // <4>
```