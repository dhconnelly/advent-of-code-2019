package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/intcode"
)

func main() {
	for _, arg := range os.Args[1:] {
		fmt.Println(arg)
		data, err := intcode.ReadProgram(arg)
		if err != nil {
			log.Fatal(err)
		}
		for _, line := range intcode.Disassemble(data) {
			fmt.Println(line)
		}
	}
}
