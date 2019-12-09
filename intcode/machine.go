package intcode

import "log"

type machine struct {
	pc      int64
	relbase int64
	data    map[int64]int64
	in      <-chan int64
	out     chan<- int64
}

func newMachine(data []int64, in <-chan int64, out chan<- int64) *machine {
	m := &machine{
		pc:      0,
		relbase: 0,
		data:    make(map[int64]int64),
		in:      in,
		out:     out,
	}
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

func (m *machine) run() {
	for ok := true; ok; {
		instr := parseInstruction(m.data[m.pc])
		if h, present := handlers[instr.op]; present {
			ok = h(m, instr)
		} else {
			log.Fatalf("bad instr at pos %d: %v", m.pc, instr)
		}
	}
	close(m.out)
}
