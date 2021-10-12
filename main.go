package main

import (
	"fmt"
	ra2 "p2/ra"
)

func main() {
	ra := ra2.New(1,"G:\\Mi unidad\\primer cuatri\\Sistemas distribuidos\\practicas\\p2\\ra\\users.txt")
	fmt.Println(ra)
	ra.PreProtocol()
}
