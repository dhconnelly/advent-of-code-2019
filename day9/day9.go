package main

import (
	"fmt"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"log"
	"os"
	"strconv"
)

func run(data []int, input int) {
	ch := make(chan int, 1)
	ch <- input
	for o := range intcode.RunProgram(data, ch) {
		fmt.Println(o)
	}
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range os.Args[2:] {
		x, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		run(data, x)
	}
}
