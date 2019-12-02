package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func readData(path string) []int {
	txt, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	toks := strings.Split(string(txt[:len(txt)-1]), ",")
	data := make([]int, len(toks))
	for i, tok := range toks {
		val, err := strconv.Atoi(tok)
		if err != nil {
			log.Fatal(err)
		}
		data[i] = val
	}
	return data
}

func execute(data []int) {
	for i := 0; i < len(data); i += 4 {
		switch data[i] {
		case 1:
			data[data[i+3]] = data[data[i+1]] + data[data[i+2]]
		case 2:
			data[data[i+3]] = data[data[i+1]] * data[data[i+2]]
		case 99:
			return
		default:
			log.Fatalf("undefined opcode at position %d: %d", i, data[i])
		}
	}
}

func executeWith(data []int, noun, verb int) int {
	local := make([]int, len(data))
	copy(local, data)
	local[1], local[2] = noun, verb
	execute(local)
	return local[0]
}

func main() {
	data := readData(os.Args[1])
	fmt.Println(executeWith(data, 12, 2))
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			if result := executeWith(data, noun, verb); result == 19690720 {
				fmt.Println(100*noun + verb)
				return
			}
		}
	}
}
