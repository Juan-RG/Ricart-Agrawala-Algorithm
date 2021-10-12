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
    "p2/ms"
    "sync"
)

type Request struct{
    Clock   int
    Pid     int
}

type Reply struct{}

type RASharedDB struct {
    OurSeqNum   int   //reloj propio
    HigSeqNum   int      //maximo reloj recibido
    OutRepCnt   int     //Cuentas de replys
    ReqCS       bool //peticion SC
    RepDefd     []int //lista de sistemas a responder
    ms          *ms.MessageSystem //Sistema de comunicación
    //canales internos
    done        chan bool //canal de finalizacion de SC
    chrep       chan bool //canal de recibir replys
    Mutex       *sync.Mutex // mutex para proteger concurrencia sobre las variables
    // TODO: completar
    Id  int // id propio
}


func New(me int, usersFile string) (*RASharedDB) {

    messageTypes := []ms.Message{Request{}, Reply{}}
    msgs := ms.New(me, usersFile, messageTypes)
    ra := RASharedDB{0, 0, 0, false, []int{}, &msgs,  make(chan bool),
                make(chan bool), &sync.Mutex{}, me}

    // TODO completar
    //llamar a pre y a post con gorutine¿?¿?

    //Arranco los procesos de recibir request y replys que esta unificado en 1   ToDo: revisar si es correcto
    go ra.receivesRequest()
    return &ra
}

/**
PROCESS WHICH INVOKES MUTUAL EXCLUSION FOR
THIS NODE
Comment Request Entry to our Critical Section;
    Variables compartidas
    P (Shared_vats)
        Comment Choose a sequence number;
        RequestingCritical_Section := TRUE;  ->Variable de acceso a la SC
        Our_Sequence_Number := Highest_Sequence_Number + l; -> Nuestra secuencia de reloj
V (Shared_vars);
Outstanding_ReplyCount := N - l;        -> Numero de replys a recibir
    //Bucle para avisar de que quiero acceder a la SC
FORj := I STEP l UNTIL N DO IFj # me THEN
    Send_Message(REQUEST(Our_Sequence_Number, me),j);
Comment sent a REQUEST message containing our sequence number and our node number to all other nodes;
Comment Now wait for a REPLY from each of the other nodes;
    //ESPERAMOS A RECIBIR LAS REPLYS
    WAITFOR (Outstanding_Reply_Count = 0);
Comment Critical Section Processing can be performed at this point;
Comment Release the Critical Section;
RequestingCritical_Section := FALSE;  -> Indicamos que ya hemos accedido a la SC
FOR j := l STEP 1 UNTIL N DO        -> Avisamos a todos los procesos de que hemos acabado
    IF Reply_Deferred[j] THEN
    BEGIN
        Reply_Deferred[j] := FALSE;
        Send_Message (REPLY, j);
    Comment send a REPLY to node j;
    END;


*/


//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol(){

    //Aumento mi reloj propio
    ra.OurSeqNum = ra.HigSeqNum + 1;
    //Indico que quiero acceder a la sc
    ra.ReqCS = true;
    //Numero de replys a recibir
    ra.OutRepCnt = len(ra.ms.Peers)         //Numero de replys a recibir

    //Aviso a todos los puntos de que quiero acceder a la seccion critica
    for id, _ := range ra.ms.Peers{
        id = id + 1                     //Indexa en id -1 por tanto hay que tratar del 1.....infi
        //si no soy yo envio un mensaje
        if id != ra.Id {
            //Rellenar                                              //Todo: Revisar los relojes segun las diapositivas
            ra.ms.Send(id, Request{ra.OurSeqNum,ra.Id})
        }
    }
    //espero a recibir todas las respuestas para entrar en la SC
    <- ra.chrep
}

//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol(){
    // TODO completar
    //cogemos el testigo para acceder a la SC
    //me falta el pasarle el fichero a cada usuario ---_____O______---
    ra.Mutex.Lock();
    //acceso de fichero
    ra.Mutex.Unlock();

    //Si alguien mas quiere acceder a la seccion critica sacamos el id del slice y le enviamos el reply
    for _, value := range ra.RepDefd {
        //Rellenar
        ra.ms.Send(value, Reply{})
    }
    //una vez enviado los reply reseteamos la lista
    ra.RepDefd = ra.RepDefd[:0]
}

func (ra *RASharedDB) Stop(){
    ra.ms.Stop()
    ra.done <- true
}
                                                                            //ToDo: Revisar actualizacion de relojs
/**
                                                                            //ToDo: Pensar como finalizar el bucle
Metodo que recibe request y replys.

*/
func (ra *RASharedDB) receivesRequest(){
    //Mirar como formatear
    //request := ra.ms.Receive();
    //request := Request{0,0}

    //Mirar como salir del bucle bien
    for {

        res := ra.ms.Receive()
        //miro el tipo de peticion
        switch v := res.(type) {                                            // ToDo: comprobar que funciona
        case Request:
            //si es request
            fmt.Println(v)
            var defer_it bool
            request := res.(Request)
            //actualizo el reloj con la >
            ra.HigSeqNum = Max(ra.HigSeqNum, request.Clock)
            //Todo: Actualizar el reloj
            //compruebo si lo ponemos en espera o si le respondemos
            defer_it = ra.ReqCS && ((request.Clock > ra.OurSeqNum) || (request.Clock > ra.OurSeqNum && request.Pid > ra.Id))
            if defer_it {
                //si espera añado el ID del proceso a la lista
                ra.RepDefd = append(ra.RepDefd, request.Pid)
            }else {
                //Todo: Revisar reply
                ra.ms.Send(ra.Id, Reply{})
            }
            break
        case Reply:
            //si recibo una reply es por que he hecho request por tanto las recibimos y descontamos del total recibidas. Para mi mejor un waitGroup --__O__--
            fmt.Println(v)
            ra.OutRepCnt--
            if ra.OutRepCnt == 0 {
                //Si es 0 aviso al proceso de que ya esta el preprotocolo
                ra.chrep <- true
            }
            break
        default:
        }
    }

}
//Metodo para escoger el numero mayor
func Max(x, y int) int {
    if x < y {
        return y
    }
    return x
}

/*
PROCESS WHICH RECEIVES REQUEST (k, j) MESSAGES
Comment k is the sequence number begin requested,
j is the node number making the request;
BOOLEAN Defer it ;
! TRUE when we cannot reply immediately
Highest_Sequence_Number :~
Maximum (Highest_Sequence_Number, k);  -> actualizamos el vector de mayor secuencia -> relojes de clase
P (Shared_vars);
// Si queremos acceder a la SC pero un proceso nos solicita acceso sera true o false segun las siguientes condiciones:
Defer it := Requesting_Critical_Section AND ((k > Our_sequence_Number) OR (k = Our_Sequence_Number ANDj > me));
V (Shared_vars);
Comment Defer_it will be TRUE if we have priority over
node j's request;
IF Defer it THEN Reply_Deferred[j] := TRUE ELSE    -> Si accedo yo a la SC aviso al cliente mas tarde
    Send_Message (REPLY, j);
    PROCESS WHICH RECEIVES REPLY MESSAGES
    Outstanding_Reply_Count := Outstanding_Reply_Count - 1;



*/