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
	"fmt"
	"p2/gestorFichero"
	"p2/ms"
	"sync"
)

type Request struct {
	Clock int
	Pid   int
	Tipo string
}

type Reply struct{}

type Token struct{
	Mensaje string
	Tipo string
}

var MATRIX = [][]bool{{false, true}, {true, true}}

type RASharedDB struct {
	OurSeqNum int               //reloj propio
	HigSeqNum int               //maximo reloj recibido
	OutRepCnt int               //Cuentas de replys
	ReqCS     bool              //peticion SC
	//RepDefd   []int             //lista de sistemas a responder
	RepDefd   []Request             //lista de sistemas a responder
	ms        *ms.MessageSystem //Sistema de comunicación
	//canales internos
	done  chan bool   //canal de finalizacion de SC
	chrep chan bool   //canal de recibir replys
	Mutex *sync.Mutex // mutex para proteger concurrencia sobre las variables
	// TODO: completar
	Id   int    // id propio
	Tipo string //tipo de accion a realizar por el proceso
	Fichero *gestorFichero.Fichero
	//	Linea string

}

func New(me int, usersFile string, tipo string, fichero *gestorFichero.Fichero) *RASharedDB {

	messageTypes := []ms.Message{Request{}, Reply{}, Token{}}
	msgs := ms.New(me, usersFile, messageTypes)
	ra := RASharedDB{0, 0, 0, false, []Request{}, &msgs, make(chan bool),
		make(chan bool), &sync.Mutex{}, me, tipo,fichero}
	//Arranco los procesos de recibir request y replys que esta unificado en 1   ToDo: revisar si es correcto
	go ra.receivesRequest()
	return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol() {
	//Variables compartidas
	ra.Mutex.Lock()
	//Aumento mi reloj propio
	ra.OurSeqNum = ra.HigSeqNum + 1

	//Indico que quiero acceder a la sc
	ra.ReqCS = true
	ra.Mutex.Unlock()

	//Numero de replys a recibir
	ra.OutRepCnt = len(ra.ms.Peers) - 1 //Numero de replys a recibir

	//Aviso a todos los puntos de que quiero acceder a la seccion critica
	for nodo, _ := range ra.ms.Peers {
		nodo = nodo + 1 //Indexa en id -1 por tanto hay que tratar del 1.....infi
		//si no soy yo envio un mensaje
		if nodo != ra.Id {
			//Rellenar                                              //Todo: Revisar los relojes segun las diapositivas
			//aumento el reloj por cada evento de envio
			ra.ms.Send(nodo, Request{ra.OurSeqNum, ra.Id, ra.Tipo})
		}
	}
	//espero a recibir todas las respuestas para entrar en la SC
	//ra.Mutex.Unlock()
	//var a bool
	_ = <-ra.chrep

}


func (ra *RASharedDB) AccesSeccionCritica(linea string) {
	//Si alguien mas quiere acceder a la seccion critica sacamos el id del slice y le enviamos el reply
	for nodo, _ := range ra.ms.Peers {
		nodo = nodo + 1 //Indexa en id -1 por tanto hay que tratar del 1.....infi
		//si no soy yo envio un mensaje
		if nodo != ra.Id {
			//Rellenar                                              //Todo: Revisar los relojes segun las diapositivas
			//aumento el reloj por cada evento de envio
			fmt.Println("envio")
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

//ToDo: Revisar actualizacion de relojs
/**
                                                                            //ToDo: Pensar como finalizar el bucle
Metodo que recibe request y replys.

*/
func (ra *RASharedDB) receivesRequest() {
	//Mirar como formatear
	//request := ra.ms.Receive();
	//request := Request{0,0}

	//Mirar como salir del bucle bien
	for {
		//ra.OurSeqNum = ra.OurSeqNum + 1
		res := ra.ms.Receive()

		//miro el tipo de peticion

		switch element := res.(type) { // ToDo: comprobar que funciona
		case Request:
			//si es request
			var defer_it bool
			//actualizo el reloj con la >
			ra.HigSeqNum = Max(ra.HigSeqNum, element.Clock)
			//Todo: Actualizar el reloj
			//compruebo si lo ponemos en espera o si le respondemos
			ra.Mutex.Lock()
			defer_it = ra.ReqCS && ((element.Clock > ra.OurSeqNum) || (element.Clock == ra.OurSeqNum && element.Pid > ra.Id)) && exclude(ra.Tipo, element.Tipo)
			ra.Mutex.Unlock()


			if defer_it {
				//si espera añado el ID del proceso a la lista
				ra.RepDefd = append(ra.RepDefd, element)
			} else {
				ra.ms.Send(element.Pid, Reply{})
			}

			break
		case Reply:
			//si recibo una reply es por que he hecho request por tanto las recibimos y descontamos del total recibidas. Para mi mejor un waitGroup --__O__--
			ra.OutRepCnt = ra.OutRepCnt - 1 //restamos un reply
			if ra.OutRepCnt == 0 {
				//Si es 0 aviso al proceso de que ya esta el preprotocolo
				ra.chrep <- true
			}
			break
		case Token:
			fmt.Println("llego")
			if element.Tipo == "escritor" {
				fmt.Println("escribo")
				ra.Fichero.EscribirFichero(element.Mensaje)
			}
			break
		}

	}

}

//              lectura escritura
//    lectura   false   true
//    escritura true    true
func exclude(opProceso string, opT string) bool {
	op_type := pasarTipoAInt(opProceso)
	op_t := pasarTipoAInt(opT)
	return MATRIX[op_type][op_t]

}

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
