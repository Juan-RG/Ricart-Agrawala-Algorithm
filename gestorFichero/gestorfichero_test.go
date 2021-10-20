package gestorFichero

import (
	"testing"
)

func TestReadFileWithLines(t *testing.T) {
	file := New("file.txt")
	if file.LeerFichero() == "" {
		t.Errorf("Error create")
	}
}

func TestWriteFileWithLines(t *testing.T) {
	file := New("file.txt")
	contenido := file.LeerFichero()
	file.EscribirFichero("hola\n")
	contenido2 := file.LeerFichero()

	if contenido == contenido2 {
		t.Errorf("Error de escritura")
	}
}
