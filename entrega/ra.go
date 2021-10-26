/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ricart-agrawala.go
* DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
 */
package ra

import (
	"p2/gestorFichero"
	"p2/ms"
	"sync"
)

type Request struct {
	Clock int
	Pid   int
	Tipo  string
}

type Reply struct{}

type Token struct {
	Mensaje string
	Tipo    string
}

//              lectura escritura
//    lectura   false   true
//    escritura true    true
var MATRIX = [][]bool{{false, true}, {true, true}}

type RASharedDB struct {
	OurSeqNum int               //reloj propio
	HigSeqNum int               //maximo reloj recibido
	OutRepCnt int               //Cuentas de replys a recibir
	ReqCS     bool              //quiero seccion critica o no
	RepDefd   []Request         //lista de sistemas a responder
	ms        *ms.MessageSystem //Sistema de comunicación

	//Canales internos y mutex
	done  chan bool   //canal de finalizacion de SC
	chrep chan bool   //canal de recibir replys
	Mutex *sync.Mutex // mutex para proteger concurrencia sobre las variables

	Id      int                    //id propio
	Tipo    string                 //tipo de accion a realizar por el proceso (escritura o lectura)
	Fichero *gestorFichero.Fichero //Puntero al fichero que se va a leer o escribir
}

//Ceramos un nuevo nodo con exclusion mutua de Ricart Agrawala
func New(me int, usersFile string, tipo string, fichero *gestorFichero.Fichero) *RASharedDB {
	messageTypes := []ms.Message{Request{}, Reply{}, Token{}}
	msgs := ms.New(me, usersFile, messageTypes)

	ra := RASharedDB{0, 0, 0, false, []Request{}, &msgs, make(chan bool),
		make(chan bool), &sync.Mutex{}, me, tipo, fichero}
	//Arranco los procesos de recibir request y replys que esta unificado en 1
	go ra.receivesMessages()

	return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol() {
	//Variables compartidas, por ello usamos mutex para asegurarnos de que solo un proceso las modifica
	ra.Mutex.Lock()
	//Aumento mi reloj propio
	ra.OurSeqNum = ra.HigSeqNum + 1

	//Indico que quiero acceder a la SC
	ra.ReqCS = true
	//Fin de modificacion de variables compartidas, liberamos el mutex
	ra.Mutex.Unlock()

	//Numero de replys a recibir
	ra.OutRepCnt = len(ra.ms.Peers) - 1 //Numero de replys a recibir

	//Aviso a todos los puntos de que quiero acceder a la seccion critica
	for nodo, _ := range ra.ms.Peers {
		nodo = nodo + 1 //Indexa en id -1 por tanto hay que tratar del 1.....infi
		//si no soy yo envio un mensaje
		if nodo != ra.Id {
			//Rellenar
			ra.ms.Send(nodo, Request{ra.OurSeqNum, ra.Id, ra.Tipo})
		}
	}

	//espero a recibir confirmacion del proceso que recibe todas las REPLY
	_ = <-ra.chrep

}

//Enviamos a todos los nodos los datos que vamos a añadir al fichero
func (ra *RASharedDB) AccesSeccionCritica(linea string) {
	//Si alguien mas quiere acceder a la seccion critica sacamos el id del slice y le enviamos el reply
	for nodo, _ := range ra.ms.Peers {
		nodo = nodo + 1 //Indexa en id -1 por tanto hay que tratar del 1.....infi
		//si no soy yo envio un mensaje
		if nodo != ra.Id {
			//Enviamos un tipo Token que contiene los datos añadidos al fichero para que el resto de nodos lo añadan
			ra.ms.Send(nodo, Token{linea, ra.Tipo})
		}
	}
}

//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol() {
	//Si alguien mas quiere acceder a la seccion critica sacamos el id del slice y le enviamos el reply
	ra.ReqCS = false

	for _, value := range ra.RepDefd {
		ra.ms.Send(value.Pid, Reply{})
	}

	//una vez enviado los reply reseteamos la lista
	ra.RepDefd = ra.RepDefd[:0]
}

func (ra *RASharedDB) Stop() {
	ra.ms.Stop()
	ra.done <- true

}

//Metodo que recibe REQUEST y REPLY.
func (ra *RASharedDB) receivesMessages() {
	for {
		res := ra.ms.Receive()

		//Comprobamos que tipo de peticion es
		switch element := res.(type) {
		case Request:
			//si es request
			var defer_it bool

			//actualizo el reloj con la >
			ra.HigSeqNum = Max(ra.HigSeqNum, element.Clock)

			//comprobamos si es necesario guardarlo en la lista de espera o podemos responderle directamente
			ra.Mutex.Lock()
			defer_it = ra.ReqCS && ((element.Clock > ra.OurSeqNum) || (element.Clock == ra.OurSeqNum && element.Pid > ra.Id)) && exclude(ra.Tipo, element.Tipo)
			ra.Mutex.Unlock()
			if defer_it {
				//si tiene que esperar ya que estamos realizando la SC, guardamos los datos para responderle posteriormente
				ra.RepDefd = append(ra.RepDefd, element)
			} else {
				ra.ms.Send(element.Pid, Reply{})
			}

			break
		case Reply:
			//si recibo una reply es por que he hecho request por tanto las recibimos y descontamos del total recibidas.
			ra.OutRepCnt = ra.OutRepCnt - 1 //restamos un reply
			if ra.OutRepCnt == 0 {
				//Si es 0 aviso al proceso de que ya esta el preprotocolo y puede comenzar la SC
				ra.chrep <- true
			}
			break
		case Token:
			//Si hemos recibido un REQUEST de un escritor, añadimos a nuestro fichero lo que escriba el escritor en el suyo
			if element.Tipo == "escritor" {
				ra.Fichero.EscribirFichero(element.Mensaje)
			}
			break
		}
	}
}

//Proceso para comprobar si los tipos de los procesos son compatibles para respoder a una REQUEST
//o es necesario posponer la respuesta
func exclude(opProceso string, opT string) bool {
	op_type := pasarTipoAInt(opProceso)
	op_t := pasarTipoAInt(opT)
	return MATRIX[op_type][op_t]

}

//Metodo para pasar el tipo de nodo a un int (para buscar de forma sencilla en la matriz)
func pasarTipoAInt(opProceso string) int {
	if opProceso == "lector" {
		return 0
	}
	return 1
}

//Metodo para escoger el numero mayor
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
