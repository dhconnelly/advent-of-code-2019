package intcode

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

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
