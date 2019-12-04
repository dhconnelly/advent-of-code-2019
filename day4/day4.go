package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func toPassword(x int) [6]byte {
	var p [6]byte
	for i := 5; i >= 0; i-- {
		p[i] = byte(x % 10)
		x /= 10
	}
	return p
}

func valid(p [6]byte) (bool, bool) {
	twoAdjacentSame, onlyTwoAdjacentSame := false, false
	matchLen := 1
	for i := 0; i < len(p)-1; i++ {
		if p[i] > p[i+1] {
			return false, false
		}
		if p[i] == p[i+1] {
			twoAdjacentSame = true
			matchLen++
		} else if matchLen == 2 {
			onlyTwoAdjacentSame = true
		} else {
			matchLen = 1
		}
	}
	return twoAdjacentSame, onlyTwoAdjacentSame || matchLen == 2
}

func countValidPasswords(from, to int) (int, int) {
	numValid1, numValid2 := 0, 0
	for i := from; i < to; i++ {
		p := toPassword(i)
		valid1, valid2 := valid(p)
		if valid1 {
			numValid1++
		}
		if valid2 {
			numValid2++
		}
	}
	return numValid1, numValid2
}

func main() {
	from, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	to, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(countValidPasswords(from, to))
}
