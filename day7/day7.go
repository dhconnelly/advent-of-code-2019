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

func execute(data []int, inputs ...int) int {
	data = copied(data)
	in, out := make(chan int, len(inputs)), make(chan int)
	for _, input := range inputs {
		in <- input
	}
	close(in)
	go run(data, in, out)
	var o int
	for o = range out {
	}
	return o
}

type seq [5]int

func availMap(avail []int) map[int]bool {
	m := make(map[int]bool)
	for _, n := range avail {
		m[n] = true
	}
	return m
}

func genSeq(
	s *seq,
	i int,
	avail map[int]bool,
	out chan<- seq,
) {
	if i >= 5 {
		out <- *s
		return
	}
	for phase, free := range avail {
		if free {
			avail[phase] = false
			s[i] = phase
			genSeq(s, i+1, avail, out)
			avail[phase] = true
		}
	}
}

func genSeqs(phases []int) chan seq {
	out := make(chan seq)
	var s seq
	go func() {
		ch := make(chan seq)
		avail := availMap(phases)
		go genSeq(&s, 0, avail, ch)
		for i := 0; i < fact(5); i++ {
			s := <-ch
			out <- s
		}
		close(out)
	}()
	return out
}

func fact(n int) int {
	if n < 1 {
		return 1
	}
	return n * fact(n-1)
}

func executeSeq(data []int, s seq) int {
	//fmt.Println("seq:", s)
	out := 0
	for _, phase := range s {
		//fmt.Printf("run %d: phase %d, input %d, ", i, phase, out)
		out = execute(data, phase, out)
		//fmt.Printf("output %d\n", out)
	}
	return out
}

func maxSignal(data []int) int {
	max := 0
	for seq := range genSeqs([]int{0, 1, 2, 3, 4}) {
		//fmt.Println(seq)
		out := executeSeq(data, seq)
		if out > max {
			max = out
		}
	}
	return max
}

func maxSignalWithLoop(data []int) int {
}

func main() {
	data := read(os.Args[1])
	fmt.Println(maxSignal(data))
}
