package intcode

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
