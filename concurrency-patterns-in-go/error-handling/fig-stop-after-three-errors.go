package main

import (
	"fmt"
	"net/http"
)

//Una buena práctica es separar en la respuesta de la gorutina la respuesta cuando todo va bien de la respuesta cuando se produce un error. En este ejemplo creamos una struct en la que se pueden informar por un lado los datos correspondientes a una respuesta exitosa, y un error

func main() {
	//Este tipo nos permite informar la respuesta existosa y el error
	type Result struct { // <1>
		Error    error          //Mensaje de error
		Response *http.Response //Respuesta existosa
	}
	//Instanciamos la función. La función admite un variant de strings, y retorna el tipo que acabamos de definir
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result { // <2>
		//Canal con la respuesta
		results := make(chan Result)

		//La gorutina
		go func() {
			defer close(results)

			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				//Empaqueta en la respuesta ambas posibilidades, el error y el exito
				result = Result{Error: err, Response: resp} // <3>
				//Implementa el for-select pattern
				select {
				case <-done:
					return
				case results <- result: // <4>
				}
			}
		}()
		return results
	}

	//Vamos a usar nuestra gorutina
	//Podemos terminar la gorutina antes de tiempo, de echo lo haremos cuando se hayan encontrado tres errores
	done := make(chan interface{})
	//En main tambien podemos hacer defer
	defer close(done)

	errCount := 0
	urls := []string{"a", "https://www.google.com", "b", "c", "d"}
	//Con range escuchamos hasta que se cierre el canal
	for result := range checkStatus(done, urls...) {
		//Si el resultado que nos llego por el canal contiene un error...
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
