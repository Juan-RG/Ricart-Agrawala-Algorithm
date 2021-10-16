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
	//ID que representara a este nodo
	var id int
	
	//ficheroNodos: Nodos con los que tendra que comunicarse para realizar la seccion critica
	//ficheroLectura: fichero del cual leera
	var ficheroNodos, ficheroLectura string
	
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

	ficheroLectura = os.Args[3]

	//Segundo de espera para volver a hacer una accion
	numeroSeg, _ := strconv.Atoi(os.Args[4])

	fmt.Println(id, ficheroNodos)

	fichero := gestorFichero.New(ficheroLectura)
	//Creamos nuevo nodo que haga uso del algoritmo de Ricart Agrawala para la exclusion mutua
	ra := ra.New(id, ficheroNodos, "lector", fichero)
	
	//Lanzamos 5 peticiones para leer el fichero
	for  i := 0; i < 5; i++ {
		//Pedimos al resto de nodos la entrada a seccion critica
		ra.PreProtocol()
		//Realizamos las operaciones necesarias en seccion critica, escribimos en el fichero
		fichero.LeerFichero()
		//Avisamos de que vamos a salir de la seccion critica
		ra.PostProtocol()
		time.Sleep(time.Second *time.Duration(numeroSeg))
	}

	fmt.Println("Terminado")
	for  {}

	fichero.CerrarDescriptor()

}
