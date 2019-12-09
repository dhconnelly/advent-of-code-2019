package intcode

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
