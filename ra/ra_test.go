package ra

import (
	"encoding/gob"
	"fmt"
	"net"
	"p2/ms"
	"reflect"
	"sync"
	"testing"
)

func TestRaNewObject(t *testing.T) {
	ra := New(1,"G:\\Mi unidad\\primer cuatri\\Sistemas distribuidos\\practicas\\p2\\ra\\users.txt")
	fmt.Println(ra.Id)
	if ra == nil{
		t.Errorf("Error create")
	}
	ra.PreProtocol()
}
type Message interface{}

func TestCastTypes(t *testing.T) {
	var wg sync.WaitGroup

	messageTypes := []ms.Message{Request{}, Reply{}}
	for _, msgTp := range messageTypes {
		gob.Register(msgTp)
	}
	wg.Add(1)
	wg.Add(1)
	go func() {
		listener, _ := net.Listen("tcp", "localhost:30000")
		for {
			select {

			default:
				conn, _ := listener.Accept()
				decoder := gob.NewDecoder(conn)
				var msg Message
				decoder.Decode(&msg)
				conn.Close()
				switch v := msg.(type) {                                            // ToDo: comprobar que funciona
				case Request:
					if reflect.TypeOf(v) != reflect.TypeOf(Request{}) {
						t.Errorf("cast error")
					}
					fmt.Println("Request: ",v)
					break
				case Reply:
					if reflect.TypeOf(v) != reflect.TypeOf(Reply{}) {
						t.Errorf("cast error")
					}
					fmt.Println("Reply: ",v)
					break
				}

			}
			wg.Done()
		}
	}()
	go func () {
		conn, _ := net.Dial("tcp", "localhost:30000")
		encoder := gob.NewEncoder(conn)
		request := Request{1,1}
		var msg Message
		msg = request
		_ = encoder.Encode(&msg)
		conn.Close()
	}()
	go func () {
		conn, _ := net.Dial("tcp", "localhost:30000")
		encoder := gob.NewEncoder(conn)
		request := Reply{}
		var msg Message
		msg = request
		_ = encoder.Encode(&msg)
		conn.Close()
	}()
	wg.Wait()
}


func TestRunRa(t *testing.T) {
	ra := New(1,"G:\\Mi unidad\\primer cuatri\\Sistemas distribuidos\\practicas\\p2\\ra\\users.txt")
	ra1 := New(2,"G:\\Mi unidad\\primer cuatri\\Sistemas distribuidos\\practicas\\p2\\ra\\users.txt")
	fmt.Println(ra.Id)
	if ra == nil{
		t.Errorf("Error create")
	}
	go ra1.PreProtocol()
	go ra.PreProtocol()
	go ra1.PostProtocol()
	go ra.PostProtocol()
	go fmt.Println(ra1)
	for  {

	}
}
/*
	p1 := New(1, "./users.txt", []Message{Request{}, Reply{}})
	p2 := New(2, "./users.txt", []Message{Request{}, Reply{}})
	p1.Send(2, Request{1})
	request := p2.Receive().(Request)

	if(request.Id != 1) {
		t.Errorf("P1 envio Request{1}, pero P2 ha recibido::Request{%d}; se esperaba Request{1}", request.Id)
	} else {
		p2.Send(1, Reply{"received"})
		msg := p1.Receive().(Reply)
		if msg.Response != "received"{
			t.Errorf("P2 envio Reply{received}, pero P1 ha recibido::Reply{%s}; se esperaba Reply{received}", msg.Response)
		}
	}



 */