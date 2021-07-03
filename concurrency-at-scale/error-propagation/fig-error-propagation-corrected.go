package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

//definimos un tipo de error. Todos los errores de nuestro modulo deberían implementar este tipo. Quien utilice nuestro modulo debe esperar, como algo normal, errores de este tipo. Si nuestro modulo arrojasé un error que no sea de este tipo se debería interpretar como un error no controlado
type MyError struct {
	//Error nativo
	Inner error
	//Mensaje de error
	Message string
	//Traza
	StackTrace string
	//data-bag
	Misc map[string]interface{}
}

//Helper para crear nuestro error
func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err, //<1>
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),        // <2>
		Misc:       make(map[string]interface{}), // <3>
	}
}

func (err MyError) Error() string {
	return err.Message
}

// "lowlevel" module

type LowLevelErr struct {
	error
}

func isGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{(wrapError(err, err.Error()))} // <1>
	}
	return info.Mode().Perm()&0100 == 0100, nil
}

// "intermediate" module

type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		return IntermediateErr{wrapError(
			err,
			"cannot run job %q: requisite binaries not available",
			id,
		)} // <1>
	} else if isExecutable == false {
		//Genera un error
		return wrapError(
			nil,
			"cannot run job %q: requisite binaries are not executable",
			id,
		)
	}

	return exec.Command(jobBinPath, "--id="+id).Run()
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := runJob("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug."
		//Si el error obtenido es del tipo que genera el modulo, será un error controlado
		if _, ok := err.(IntermediateErr); ok {
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}
