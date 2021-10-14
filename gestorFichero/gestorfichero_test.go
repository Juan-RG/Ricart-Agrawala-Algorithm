package gestorFichero

import (
	"fmt"
	"testing"
)

func TestReadFileWithLines(t *testing.T) {
	file := new("file.txt")
	if file.LeerFichero() == "" {
		t.Errorf("Error create")
	}
}

func TestWriteFileWithLines(t *testing.T) {
	file := new("file.txt")
	contenido := file.LeerFichero()
	file.EscribirFichero("hola\n")
	contenido2 := file.LeerFichero()

	fmt.Println("--- ", contenido)
	fmt.Println("--- ", contenido2)
	if contenido == contenido2 {
		t.Errorf("Error de escritura")
	}
}
