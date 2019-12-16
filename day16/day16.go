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

func copied(data []int) []int {
	data2 := make([]int, len(data))
	copy(data2, data)
	return data2
}

func fft(signal []int, phases int) []int {
	signal = copied(signal)
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

func indexRepeated(signal []int, i int) int {
	return signal[i%len(signal)]
}

func sliceSignal(signal []int, from, to int) []int {
	sliced := make([]int, to-from)
	for i := range sliced {
		sliced[i] = indexRepeated(signal, from+i)
	}
	return sliced
}

func extractMessage(signal []int, reps, phases, offset, digits int) []int {
	msg := sliceSignal(signal, offset, len(signal)*reps)
	n := len(msg)

	fmt.Println("findings coefs")
	coefs := make([][]int, n)
	for i := 0; i < n; i++ {
		coefs[i] = make([]int, n-i)
		for j := i; j < n; j++ {
			coefs[i][j-i] = coef(offset+i+1, offset+j+1)
		}
	}

	scratch := make([]int, n)
	for ; phases > 0; phases-- {
		fmt.Println("phase:", phases)
		for i := 0; i < n; i++ {
			sum := 0
			for j := i; j < n; j++ {
				sum += (coefs[i][j-i] * msg[j])
			}
			scratch[i] = ints.Abs(sum) % 10
		}
		msg, scratch = scratch, msg
	}

	return msg[:digits]
}

func main() {
	txt, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	v := parseVec(bytes.TrimSpace(txt))
	i := offset(txt)

	const digits = 8
	const phases = 100
	const reps = 10000

	//fmt.Println(toSignal(fft(v, phases)[:digits]))
	fmt.Println("offset:", i)
	fmt.Println(toSignal(extractMessage(v, reps, phases, i, digits)))
}
