package main

import (
	"encoding/gob"
	"fmt"
	"github.com/DistributedClocks/GoVector/govec"
	"net"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	Logger := govec.InitGoVector("ServerLogs", "FileServerLog", govec.GetDefaultConfig())
	listener, err := net.Listen("tcp", "localhost:8081")
	checkError(err)
	for {
		conn, err := listener.Accept()
		checkError(err)
		decoder := gob.NewDecoder(conn)
		var buf []byte
		err = decoder.Decode(&buf)
		fmt.Println("alo")
		fmt.Println(string(buf))
		checkError(err)
		var idProces int
		opts := govec.GetDefaultLogOptions()
		Logger.UnpackReceive("Received Message From Client", buf[0:], &idProces, opts)
		conn.Close()
	}
}
