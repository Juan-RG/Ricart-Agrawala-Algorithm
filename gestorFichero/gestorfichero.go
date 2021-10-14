package gestorFichero

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type fichero struct {
	f *os.File
}

func new(nombre string) *fichero {
	f, err := os.OpenFile(nombre, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fichero := fichero{f}
	return &fichero
}

func (file *fichero) LeerFichero() string {
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

func (file *fichero) EscribirFichero(fragmento string) {
	var data []byte

	// Leer el fichero
	data, err := ioutil.ReadAll(file.f)
	if err != nil {
		fmt.Printf("Error leyendo fichero: %s\n", err)
		os.Exit(1)
	}

	// Agrego contenido
	data = append(data, []byte(fragmento)...)

	// Guardar contenido
	_, err = file.f.Write(data)
	if err != nil {
		fmt.Println(err)
	}

}

func (file *fichero) CerrarDescriptor() {
	file.f.Close()
}
