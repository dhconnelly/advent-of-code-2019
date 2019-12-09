package intcode

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

type handler func(m *machine, instr instruction) bool

var handlers = map[opcode]handler{
	add: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		s := m.data[m.pc+3]
		m.set(s, l+r, instr.modes[2])
		m.pc += instr.arity + 1
		return true
	},

	mul: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		s := m.data[m.pc+3]
		m.set(s, l*r, instr.modes[2])
		m.pc += instr.arity + 1
		return true
	},

	read: func(m *machine, instr instruction) bool {
		s := m.data[m.pc+1]
		m.set(s, m.read(), instr.modes[0])
		m.pc += instr.arity + 1
		return true
	},

	print: func(m *machine, instr instruction) bool {
		v := m.get(m.pc+1, instr.modes[0])
		m.write(v)
		m.pc += instr.arity + 1
		return true
	},

	jmpif: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		if l != 0 {
			m.pc = r
		} else {
			m.pc += instr.arity + 1
		}
		return true
	},

	jmpnot: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		if l == 0 {
			m.pc = r
		} else {
			m.pc += instr.arity + 1
		}
		return true
	},

	lt: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		s := m.data[m.pc+3]
		if l < r {
			m.set(s, 1, instr.modes[2])
		} else {
			m.set(s, 0, instr.modes[2])
		}
		m.pc += instr.arity + 1
		return true
	},

	eq: func(m *machine, instr instruction) bool {
		l := m.get(m.pc+1, instr.modes[0])
		r := m.get(m.pc+2, instr.modes[1])
		s := m.data[m.pc+3]
		if l == r {
			m.set(s, 1, instr.modes[2])
		} else {
			m.set(s, 0, instr.modes[2])
		}
		m.pc += instr.arity + 1
		return true
	},

	adjrel: func(m *machine, instr instruction) bool {
		v := m.get(m.pc+1, instr.modes[0])
		m.relbase += v
		m.pc += instr.arity + 1
		return true
	},

	halt: func(m *machine, instr instruction) bool {
		return false
	},
}
