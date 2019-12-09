package main

import (
	"fmt"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"log"
	"os"
)

func run(data []int, input int) int {
	ch := make(chan int, 1)
	ch <- input
	return <-intcode.RunProgram(data, ch)
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(run(data, 1))
	fmt.Println(run(data, 2))
}
