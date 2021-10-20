package main

import (
	"encoding/gob"
	"fmt"
	"github.com/DistributedClocks/GoVector/govec"
	"net"
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


	nombreFileLog := "lector" + strconv.Itoa(id)
	Logger := govec.InitGoVector(nombreFileLog, nombreFileLog, govec.GetDefaultConfig())
	fichero := gestorFichero.New(ficheroLectura)
	opts := govec.GetDefaultLogOptions()
	//Creamos nuevo nodo que haga uso del algoritmo de Ricart Agrawala para la exclusion mutua
	ra := ra.New(id, ficheroNodos, "lector", fichero)
	time.Sleep(time.Second * time.Duration(5))
	//Lanzamos 5 peticiones para leer el fichero
	for i := 0; i < 5; i++ {
		outBuf := Logger.PrepareSend("Acceder SC", id, opts)
		//Pedimos al resto de nodos la entrada a seccion critica
		ra.PreProtocol()
		//Realizamos las operaciones necesarias en seccion critica, escribimos en el fichero
		fichero.LeerFichero()
		enviarServerLogs(outBuf)

		//Avisamos de que vamos a salir de la seccion critica
		ra.PostProtocol()
		time.Sleep(time.Second * time.Duration(numeroSeg))

	}

	fmt.Println("Terminado")
	for {
	}

	fichero.CerrarDescriptor()

}

func enviarServerLogs(buf []byte) {
	conn, _ := net.Dial("tcp", "localhost:8081")
	encoder := gob.NewEncoder(conn)
	_ = encoder.Encode(&buf)
	defer conn.Close()
}
