package main

import "fmt"

func main() {
	RepDefd := []int{9, 8, 7, 6}

	for i := len(RepDefd)-1; i >= 0; i-- {
		fmt.Println(RepDefd[i])
	}

}

