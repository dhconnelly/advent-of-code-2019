package intcode

type mode int

const (
	pos mode = iota
	imm
	rel
)

type opcode int

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
