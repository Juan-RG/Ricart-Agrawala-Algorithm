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
	var ficheroNodos string
	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("error cast id nodo")
		id = 1
	}

	if len(os.Args) > 2 && os.Args[2] != "" {
		ficheroNodos = os.Args[2]
	} else {
		ficheroNodos = "G:\\Mi unidad\\primer cuatri\\Sistemas distribuidos\\practicas\\p2\\users.txt"
	}
	fmt.Println(id, ficheroNodos)
	ra := ra.New(id, ficheroNodos, "escritor")

	for {
		fmt.Println(id, "preproto")
		ra.PreProtocol()
		fmt.Println(id, "leerFichero")
		gestorFichero.EscribirFichero("hola-------------------------")
		fmt.Println(id, "postproto")
		ra.PostProtocol()
		fmt.Println(id, "salgo")
		time.Sleep(time.Second * 3)
	}

}
