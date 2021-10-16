package gestorFichero

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Fichero struct {
	f *os.File
	Mutex *sync.Mutex // mutex para proteger concurrencia sobre las variables
}

func New(nombre string) *Fichero {
	f, err := os.OpenFile(nombre, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fichero := Fichero{f,&sync.Mutex{}}
	return &fichero
}

func (file *Fichero) LeerFichero() string {
	var data []byte
	file.Mutex.Lock()
	// Leer el fichero
	data, err := ioutil.ReadAll(file.f)
	if err != nil {
		fmt.Printf("Error leyendo fichero: %s\n", err)
		os.Exit(1)
	}
	datosComoString := string(data)
	file.Mutex.Unlock()
	return datosComoString
}

func (file *Fichero) EscribirFichero(fragmento string) {
	var data []byte

	file.Mutex.Lock()
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
	file.Mutex.Unlock()

}

func (file *Fichero) CerrarDescriptor() {
	file.f.Close()
}
