package gestorFichero

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Fichero struct {
	f *os.File
	//Mutex *sync.Mutex // mutex para proteger concurrencia sobre las variables
}

//Devolvemos el puntero del fichero que recibimos
func New(nombre string) *Fichero {
	f, err := os.OpenFile(nombre, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//fichero := Fichero{f,&sync.Mutex{}}
	fichero := Fichero{f}

	return &fichero
}

//Devuelve el contenido del fichero
func (file *Fichero) LeerFichero() string {
	var data []byte
	// Leer el fichero
	data, err := ioutil.ReadAll(file.f)
	if err != nil {
		fmt.Printf("Error leyendo fichero: %s\n", err)
		os.Exit(1)
	}
	datosComoString := string(data)
	return datosComoString
}

//Añade al contenido del fichero, lo que recibe del parametro fragmento
func (file *Fichero) EscribirFichero(fragmento string) {
	var data []byte

	// Leer el fichero
	data, err := ioutil.ReadAll(file.f)
	if err != nil {
		fmt.Printf("Error leyendo fichero: %s\n", err)
		os.Exit(1)
	}

	// Agrego contenido
	data = append(data, []byte(fragmento)...)
	data = append(data, []byte("\n")...)

	// Guardar contenido
	_, err = file.f.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}

//Cerramos el descriptor del fichero
func (file *Fichero) CerrarDescriptor() {
	file.f.Close()
}
