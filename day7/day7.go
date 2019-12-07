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
	1: 3, 2: 3, 3: 1, 4: 1, 5: 2, 6: 2, 7: 3, 8: 3, 99: 0,
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

func atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func read(path string) []int {
	txt, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	toks := strings.Split(strings.TrimSpace(string(txt)), ",")
	data := make([]int, len(toks))
	for i, tok := range toks {
		data[i] = atoi(tok)
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
			i += instr.params + 1
		case 2:
			l := get(data, i+1, instr.modes[0])
			r := get(data, i+2, instr.modes[1])
			s := data[i+3]
			data[s] = l * r
			i += instr.params + 1
		case 3:
			s := data[i+1]
			data[s] = <-in
			i += instr.params + 1
		case 4:
			v := get(data, i+1, instr.modes[0])
			out <- v
			i += instr.params + 1
		case 5:
			l := get(data, i+1, instr.modes[0])
			r := get(data, i+2, instr.modes[1])
			if l != 0 {
				i = r
			} else {
				i += instr.params + 1
			}
		case 6:
			l := get(data, i+1, instr.modes[0])
			r := get(data, i+2, instr.modes[1])
			if l == 0 {
				i = r
			} else {
				i += instr.params + 1
			}
		case 7:
			l := get(data, i+1, instr.modes[0])
			r := get(data, i+2, instr.modes[1])
			s := data[i+3]
			if l < r {
				data[s] = 1
			} else {
				data[s] = 0
			}
			i += instr.params + 1
		case 8:
			l := get(data, i+1, instr.modes[0])
			r := get(data, i+2, instr.modes[1])
			s := data[i+3]
			if l == r {
				data[s] = 1
			} else {
				data[s] = 0
			}
			i += instr.params + 1
		case 99:
			close(out)
			return
		}
	}
}

func copied(data []int) []int {
	data2 := make([]int, len(data))
	copy(data2, data)
	return data2
}

func execute(data []int, phase int, signals <-chan int) chan int {
	data = copied(data)
	in, out := make(chan int), make(chan int)
	go func() {
		in <- phase
		for signal := range signals {
			in <- signal
		}
		close(in)
	}()
	go run(data, in, out)
	return out
}

type seq [5]int

func genSeqsRec(nums []int, used map[int]bool, until int) []seq {
	if until == 0 {
		return []seq{{}}
	}
	var seqs []seq
	for _, num := range nums {
		if !used[num] {
			used[num] = true
			for _, recSeq := range genSeqsRec(nums, used, until-1) {
				recSeq[until-1] = num
				seqs = append(seqs, recSeq)
			}
			used[num] = false
		}
	}
	return seqs
}

func genSeqs(nums []int) []seq {
	return genSeqsRec(nums, make(map[int]bool), len(nums))
}

func executeSeq(data []int, s seq) int {
	out := 0
	for _, phase := range s {
		in := make(chan int, 1)
		in <- out
		close(in)
		out = <-execute(data, phase, in)
	}
	return out
}

func executeWithFeedback(data []int, s seq) int {
	// send the initial input
	in := make(chan int, 1)
	in <- 0

	// pipe the amplifiers together
	out := in
	for _, phase := range s {
		out = execute(data, phase, out)
	}

	// pipe the output back into the input, but keep track of it
	var o int
	for o = range out {
		in <- o
	}

	// last output is the output signal
	return o
}

func maxSignal(data []int, exec func([]int, seq) int, nums []int) int {
	max := 0
	for _, seq := range genSeqs(nums) {
		out := exec(data, seq)
		if out > max {
			max = out
		}
	}
	return max
}

func main() {
	data := read(os.Args[1])
	fmt.Println(maxSignal(data, executeSeq, []int{0, 1, 2, 3, 4}))
	fmt.Println(maxSignal(data, executeWithFeedback, []int{5, 6, 7, 8, 9}))
}
