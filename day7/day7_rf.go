package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/intcode"
)

func execute(data []int, phase int, signals <-chan int) <-chan int {
	in := make(chan int)
	go func() {
		in <- phase
		for signal := range signals {
			in <- signal
		}
		close(in)
	}()
	return intcode.RunProgram(data, in)
}

type seq [5]int

func genSeqsRec(nums []int, used map[int]bool, until int) []seq {
	if until == 0 {
		return []seq{{}}
	}
	var seqs []seq
	for _, num := range nums {
		if !used[num] {
			used[num] = true
			for _, recSeq := range genSeqsRec(nums, used, until-1) {
				recSeq[until-1] = num
				seqs = append(seqs, recSeq)
			}
			used[num] = false
		}
	}
	return seqs
}

func genSeqs(nums []int) []seq {
	return genSeqsRec(nums, make(map[int]bool), len(nums))
}

func executeSeq(data []int, s seq) int {
	out := 0
	for _, phase := range s {
		in := make(chan int, 1)
		in <- out
		close(in)
		out = <-execute(data, phase, in)
	}
	return out
}

func executeWithFeedback(data []int, s seq) int {
	// send the initial input
	in := make(chan int, 1)
	in <- 0

	// pipe the amplifiers together
	var out <-chan int = in
	for _, phase := range s {
		out = execute(data, phase, out)
	}

	// pipe the output back into the input, but keep track of it
	var o int
	for o = range out {
		in <- o
	}

	// last output is the output signal
	return o
}

func maxSignal(data []int, exec func([]int, seq) int, nums []int) int {
	max := 0
	for _, seq := range genSeqs(nums) {
		out := exec(data, seq)
		if out > max {
			max = out
		}
	}
	return max
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(maxSignal(data, executeSeq, []int{0, 1, 2, 3, 4}))
	fmt.Println(maxSignal(data, executeWithFeedback, []int{5, 6, 7, 8, 9}))
}
