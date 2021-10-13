package main

import (
	"io/ioutil"
	"log"
	"p2/ra"
	"sync"
)
func LeerFichero() string{
	datosComoBytes, err := ioutil.ReadFile("./datos.txt")
	if err != nil {
		log.Fatal(err)
	}
	// convertir el arreglo a string
	datosComoString := string(datosComoBytes)

	return datosComoString
}
func EscribirFichero(fragmento string){
	// write the whole body at once
	err := ioutil.WriteFile("./datos.txt", []byte(fragmento), 0644)
	if err != nil {
		panic(err)
	}
}
func main() {
	var wg sync.WaitGroup
	ra1 := ra.New(1,"G:\\Mi unidad\\primer cuatri\\Sistemas distribuidos\\practicas\\p2\\ra\\users.txt")
	ra2 := ra.New(2,"G:\\Mi unidad\\primer cuatri\\Sistemas distribuidos\\practicas\\p2\\ra\\users.txt")
	wg.Add(1)
	go ra2.PreProtocol()

	ra1.PreProtocol()

	ra1.PostProtocol()
	wg.Wait()


}
