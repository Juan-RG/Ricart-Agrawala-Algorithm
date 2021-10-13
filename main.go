package main

import (
	"io/ioutil"
	"log"
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

	//ra := ra.New(1,"G:\\Mi unidad\\primer cuatri\\Sistemas distribuidos\\practicas\\p2\\ra\\users.txt")
	//fmt.Println(ra)

}
