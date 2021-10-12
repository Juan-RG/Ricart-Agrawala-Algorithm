package main

import "fmt"

func main() {
/*
	ra := ra2.New(1,"G:\\Mi unidad\\primer cuatri\\Sistemas distribuidos\\practicas\\p2\\ra\\users.txt")
	fmt.Println(ra)
	ra.PreProtocol()
*/
	var s []int
	printSlice(s)

	// append works on nil slices.
	s = append(s, 1)
	printSlice(s)

	// The slice grows as needed.
	s = append(s, 1)
	printSlice(s)

	// We can add more than one element at a time.
	s = append(s, 1, 1, 1)
	printSlice(s)
	for i, value := range s {
		fmt.Println("valores i: ", i , " value: ", value)
	}
	s = s[:0]
	printSlice(s)
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
