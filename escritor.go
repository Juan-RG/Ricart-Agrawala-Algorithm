package main

import (
	"fmt"
	"os"
	"p2/gestorFichero"
	"p2/ra"
	"strconv"
	"time"
)

func main() {
	var id int
	var ficheroNodos, ficheroEscritura, linea string
	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("error cast id nodo")
		id = 1
	}

	if os.Args[2] != "" {
		ficheroNodos = os.Args[2]
	} else {
		ficheroNodos = "G:\\Mi unidad\\primer cuatri\\Sistemas distribuidos\\practicas\\p2\\users.txt"
	}
	ficheroEscritura = os.Args[3]
	linea = os.Args[4]
	numeroSeg, _ := strconv.Atoi(os.Args[5])

	fmt.Println(id, ficheroNodos)
	fichero := gestorFichero.New(ficheroEscritura)
	ra := ra.New(id, ficheroNodos, "escritor", fichero)
	for  i := 0; i < 5; i++ {
		ra.PreProtocol()
		fichero.EscribirFichero(linea)
		ra.AccesSeccionCritica(linea)
		ra.PostProtocol()
		time.Sleep(time.Second *time.Duration(numeroSeg))

	}
	fmt.Println("final")

	for  {
		
	}
	fichero.CerrarDescriptor()

}
