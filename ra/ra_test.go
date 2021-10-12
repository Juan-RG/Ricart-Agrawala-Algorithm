package ra

import (
	"fmt"
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