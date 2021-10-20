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
	/*	Logger := govec.InitGoVector("client", "clientlogfile", govec.GetDefaultConfig())


	 */

	//ID que representara a este nodo
	var id int

	//ficheroNodos: Nodos con los que tendra que comunicarse para realizar la seccion critica
	//ficheroEscritura: fichero que tendra cada nodo para escribir
	//linea: nuevos datos a añadir al fichero
	var ficheroNodos, ficheroEscritura, linea string
	fmt.Println(linea, ficheroNodos)
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

	//Segundo de espera para volver a hacer una accion
	numeroSeg, _ := strconv.Atoi(os.Args[5])

	//logger del escritor
	var nombreFileLog string
	nombreFileLog = "escritor" + strconv.Itoa(id)
	fmt.Println(nombreFileLog)
	Logger := govec.InitGoVector(nombreFileLog, nombreFileLog, govec.GetDefaultConfig())
	opts := govec.GetDefaultLogOptions()
	fichero := gestorFichero.New(ficheroEscritura)
	//Creamos nuevo nodo que haga uso del algoritmo de Ricart Agrawala para la exclusion mutua
	ra := ra.New(id, ficheroNodos, "escritor", fichero)

	//Lanzamos 5 peticiones de seccion critica para escribir en el fichero
	for i := 0; i < 500; i++ {

		outBuf := Logger.PrepareSend("Acceder SC", id, opts)
		//Pedimos al resto de nodos la entrada a seccion critica
		ra.PreProtocol()
		//Realizamos las operaciones necesarias en seccion critica, escribimos en el fichero
		fichero.EscribirFichero(linea)
		enviarServerLogs(outBuf)
		//
		ra.AccesSeccionCritica(linea)
		//Avisamos de que vamos a salir de la seccion critica
		ra.PostProtocol()
		time.Sleep(time.Second * time.Duration(numeroSeg))
	}

	for {

	}
	//ra.Stop()

	defer fichero.CerrarDescriptor()

}

func enviarServerLogs(buf []byte) {
	conn, _ := net.Dial("tcp", "localhost:8081")
	encoder := gob.NewEncoder(conn)
	_ = encoder.Encode(&buf)
	defer conn.Close()
}
