package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/dhconnelly/advent-of-code-2019/ints"
)

func parseVec(b []byte) []int {
	v := make([]int, len(b))
	for i, a := range b {
		v[i] = int(a - '0')
	}
	return v
}

func coef(row, col int) int {
	if col < row {
		return 0
	}
	switch (col / row) % 4 {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 0
	default:
		return -1
	}
}

func fft(signal []int, phases int) []int {
	scratch := make([]int, len(signal))
	for ; phases > 0; phases-- {
		for i := 0; i < len(signal); i++ {
			sum := 0
			for j := 0; j < len(signal); j++ {
				sum += coef(i+1, j+1) * signal[j]
			}
			scratch[i] = ints.Abs(sum) % 10
		}
		signal, scratch = scratch, signal
	}
	return signal
}

func toSignal(signal []int) string {
	b := make([]byte, len(signal))
	for i, v := range signal {
		b[i] = byte(v + '0')
	}
	return string(b)
}

func offset(b []byte) int {
	i, err := strconv.Atoi(string(b[:7]))
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func main() {
	txt, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	v := parseVec(bytes.TrimSpace(txt))
	const digits = 8
	fmt.Println(toSignal(fft(v, 100)[:digits]))
	fmt.Println(offset(txt))
}
