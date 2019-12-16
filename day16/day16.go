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

func coefCol(col, width int) []int {
	v := make([]int, width)
	for i := 0; i < width; i++ {
		v[i] = coef(i+1, col)
	}
	return v
}

func coefColPow(col, width, n int) []int {
	v := coefCol(col, width)
	for n--; n > 0; n-- {
		v2 := make([]int, len(v))
		for i := 0; i < width; i++ {
			sum := 0
			for j := 0; j < width; j++ {
				sum += coef(i+1, j+1) * v[j]
			}
			v2[i] = sum
		}
		v = v2
	}
	return v
}

func transpose(mat [][]int) {
	for i := 0; i < len(mat); i++ {
		for j := i + 1; j < len(mat); j++ {
			mat[i][j], mat[j][i] = mat[j][i], mat[i][j]
		}
	}
}

func coefMatPow(width, n int) [][]int {
	mat := make([][]int, width)
	for j := 0; j < width; j++ {
		mat[j] = coefColPow(j+1, width, n)
	}
	transpose(mat)
	return mat
}

func coefMat(width int) [][]int {
	mat := make([][]int, width)
	for i := 0; i < width; i++ {
		mat[i] = make([]int, width)
		for j := 0; j < width; j++ {
			mat[i][j] = coef(i+1, j+1)
		}
	}
	return mat
}

func printMat(mat [][]int) {
	for _, row := range mat {
		for _, col := range row {
			fmt.Printf("%3d", col)
		}
		fmt.Println()
	}
}

func dot(x, y []int) int {
	sum := 0
	for i := 0; i < len(x); i++ {
		sum += x[i] * y[i]
	}
	return ints.Abs(sum) % 10
}

func fft(signal []int, phases int) []int {
	mat := coefMat(len(signal))
	signal2 := make([]int, len(signal))
	for ; phases > 0; phases-- {
		copy(signal2, signal)
		for i := 0; i < len(signal); i++ {
			signal2[i] = dot(mat[i], signal)
		}
		copy(signal, signal2)
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
