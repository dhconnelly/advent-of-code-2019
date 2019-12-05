package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type mode int

const (
	POS mode = iota
	IMM
)

type instr struct {
	opcode int
	params int
	modes  []mode
}

var opcodeToParam = map[int]int{
	1: 3, 2: 3, 3: 1, 4: 1, 99: 0,
}

func parseInstr(i int) instr {
	var in instr
	in.opcode = i % 100
	in.params = opcodeToParam[in.opcode]
	for i /= 100; len(in.modes) < in.params; i /= 10 {
		in.modes = append(in.modes, mode(i%10))
	}
	return in
}

func read(path string) []int {
	txt, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	toks := strings.Split(strings.TrimSpace(string(txt)), ",")
	data := make([]int, len(toks))
	for i, tok := range toks {
		data[i], err = strconv.Atoi(tok)
		if err != nil {
			log.Fatal(err)
		}
	}
	return data
}

func get(data []int, i int, m mode) int {
	v := data[i]
	switch m {
	case POS:
		return data[v]
	case IMM:
		return v
	}
	log.Fatalf("unknown mode: %d", m)
	return 0
}

func run(data []int, in <-chan int, out chan<- int) {
	for i := 0; i < len(data); {
		instr := parseInstr(data[i])
		switch instr.opcode {
		case 1:
			l := get(data, i+1, instr.modes[0])
			r := get(data, i+2, instr.modes[1])
			s := data[i+3]
			data[s] = l + r
		case 2:
			l := get(data, i+1, instr.modes[0])
			r := get(data, i+2, instr.modes[1])
			s := data[i+3]
			data[s] = l * r
		case 3:
			s := data[i+1]
			data[s] = <-in
		case 4:
			v := get(data, i+1, instr.modes[0])
			out <- v
		case 99:
			close(out)
			break
		}
		i += instr.params + 1
	}
}

func main() {
	data := read(os.Args[1])
	in, out := make(chan int, 1), make(chan int)
	in <- 1
	go run(data, in, out)
	for o := range out {
		fmt.Println(o)
	}
	close(in)
}
