package intcode

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Opcode int

const (
	ADD    Opcode = 1
	MUL           = 2
	READ          = 3
	PRINT         = 4
	JMPIF         = 5
	JMPNOT        = 6
	LT            = 7
	EQ            = 8
	HALT          = 99
)

func ReadProgram(path string) ([]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var data []int
	for r := bufio.NewReader(f); ; {
		tok, err := r.ReadString(',')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("bad int in %s: %w", path, err)
		}
		if len(tok) > 0 {
			var i int
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
