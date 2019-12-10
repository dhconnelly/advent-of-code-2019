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

// Retrieves a value according to the specified mode.
//
// * In immediate mode, returns the value stored at the given address.
//
// * In position mode, the value stored at the address is interpreted
//   as a *pointer* to the value that should be returned.
//
// * In relative mode, the machine's current relative base is interpreted
//   as a pointer, and the value stored at the address is interpreted
//   as an offset to that pointer. The value stored at the *resulting*
//   address is returned.
//
func (m *machine) get(addr int64, md mode) int64 {
	v := m.data[addr]
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

// Sets a value according to the specified mode.
//
// * In position mode, the value stored at the given address specifies
//   the address to which the value should be written.
//
// * In relative mode, the value stored at the given address specifies
//   an offset to the relative base, and the sum of the offset and the
//   base specifies the address to which the value should be written.
//
func (m *machine) set(addr, val int64, md mode) {
	v := m.data[addr]
	switch md {
	case pos:
		m.data[v] = val
	case rel:
		m.data[v+m.relbase] = val
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
