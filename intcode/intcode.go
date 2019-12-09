package intcode

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type mode int

const (
	pos mode = iota
	imm
	rel
)

type opcode int

const (
	add    opcode = 1
	mul           = 2
	read          = 3
	print         = 4
	jmpif         = 5
	jmpnot        = 6
	lt            = 7
	eq            = 8
	adjrel        = 9
	halt          = 99
)

var opcodeToArity = map[opcode]int64{
	add:    3,
	mul:    3,
	read:   1,
	print:  1,
	jmpif:  2,
	jmpnot: 2,
	lt:     3,
	eq:     3,
	adjrel: 1,
	halt:   0,
}

type instruction struct {
	op    opcode
	arity int64
	modes []mode
}

func parseInstruction(i int64) instruction {
	var in instruction
	in.op = opcode(i % 100)
	in.arity = opcodeToArity[in.op]
	for i /= 100; int64(len(in.modes)) < in.arity; i /= 10 {
		in.modes = append(in.modes, mode(i%10))
	}
	return in
}

type machine struct {
	relbase int64
	data    map[int64]int64
	in      <-chan int64
	out     chan<- int64
}

func newMachine(data []int64, in <-chan int64, out chan<- int64) *machine {
	m := &machine{0, make(map[int64]int64), in, out}
	for i, v := range data {
		m.data[int64(i)] = v
	}
	return m
}

func (m *machine) get(i int64, md mode) int64 {
	v := m.data[i]
	switch md {
	case pos:
		return m.data[v]
	case imm:
		return v
	case rel:
		return m.data[v+m.relbase]
	}
	log.Fatalf("unknown mode: %d", md)
	return 0
}

func (m *machine) set(i, x int64, md mode) {
	switch md {
	case pos:
		m.data[i] = x
	case rel:
		m.data[i+m.relbase] = x
	default:
		log.Fatalf("bad mode for write: %d", md)
	}
}

func (m *machine) read() int64 {
	return <-m.in
}

func (m *machine) write(x int64) {
	m.out <- x
}

type handler func(m *machine, pc int64, instr instruction) (int64, bool)

var handlers = map[opcode]handler{
	add: func(m *machine, pc int64, instr instruction) (int64, bool) {
		l := m.get(pc+1, instr.modes[0])
		r := m.get(pc+2, instr.modes[1])
		s := m.data[pc+3]
		m.set(s, l+r, instr.modes[2])
		return pc + instr.arity + 1, true
	},

	mul: func(m *machine, pc int64, instr instruction) (int64, bool) {
		l := m.get(pc+1, instr.modes[0])
		r := m.get(pc+2, instr.modes[1])
		s := m.data[pc+3]
		m.set(s, l*r, instr.modes[2])
		return pc + instr.arity + 1, true
	},

	read: func(m *machine, pc int64, instr instruction) (int64, bool) {
		s := m.data[pc+1]
		m.set(s, m.read(), instr.modes[0])
		return pc + instr.arity + 1, true
	},

	print: func(m *machine, pc int64, instr instruction) (int64, bool) {
		v := m.get(pc+1, instr.modes[0])
		m.write(v)
		return pc + instr.arity + 1, true
	},

	jmpif: func(m *machine, pc int64, instr instruction) (int64, bool) {
		l := m.get(pc+1, instr.modes[0])
		r := m.get(pc+2, instr.modes[1])
		if l != 0 {
			return r, true
		} else {
			return pc + instr.arity + 1, true
		}
	},

	jmpnot: func(m *machine, pc int64, instr instruction) (int64, bool) {
		l := m.get(pc+1, instr.modes[0])
		r := m.get(pc+2, instr.modes[1])
		if l == 0 {
			return r, true
		} else {
			return pc + instr.arity + 1, true
		}
	},

	lt: func(m *machine, pc int64, instr instruction) (int64, bool) {
		l := m.get(pc+1, instr.modes[0])
		r := m.get(pc+2, instr.modes[1])
		s := m.data[pc+3]
		if l < r {
			m.set(s, 1, instr.modes[2])
		} else {
			m.set(s, 0, instr.modes[2])
		}
		return pc + instr.arity + 1, true
	},

	eq: func(m *machine, pc int64, instr instruction) (int64, bool) {
		l := m.get(pc+1, instr.modes[0])
		r := m.get(pc+2, instr.modes[1])
		s := m.data[pc+3]
		if l == r {
			m.set(s, 1, instr.modes[2])
		} else {
			m.set(s, 0, instr.modes[2])
		}
		return pc + instr.arity + 1, true
	},

	adjrel: func(m *machine, pc int64, instr instruction) (int64, bool) {
		v := m.get(pc+1, instr.modes[0])
		m.relbase += v
		return pc + instr.arity + 1, true
	},

	halt: func(m *machine, pc int64, instr instruction) (int64, bool) {
		return 0, false
	},
}

func (m *machine) run() {
	for pc, ok := int64(0), true; ok; {
		instr := parseInstruction(m.data[pc])
		if h, present := handlers[instr.op]; present {
			pc, ok = h(m, pc, instr)
		} else {
			log.Fatalf("bad instr at pos %d: %v", pc, instr)
		}
		if !ok {
			close(m.out)
		}
	}
}

func RunProgram(data []int64, in <-chan int64) <-chan int64 {
	out := make(chan int64)
	m := newMachine(data, in, out)
	go m.run()
	return out
}

func ReadProgram(path string) ([]int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var data []int64
	for r := bufio.NewReader(f); ; {
		tok, err := r.ReadString(',')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("bad int in %s: %w", path, err)
		}
		if len(tok) > 0 {
			var i int64
			if _, err := fmt.Sscanf(tok, "%d", &i); err != nil {
				return nil, fmt.Errorf("bad int in %s: %w", path, err)
			}
			data = append(data, i)
		}
		if err == io.EOF {
			break
		}
	}
	return data, nil
}
