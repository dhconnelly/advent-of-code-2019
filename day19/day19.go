package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"github.com/dhconnelly/advent-of-code-2019/ints"
)

type state int64

const (
	stationary = 0
	pulled     = 1
)

type drone struct {
	prog []int64
}

func newDrone(prog []int64) *drone {
	return &drone{ints.Copied64(prog)}
}

func (d *drone) test(x, y int, debug bool) state {
	prog := ints.Copied64(d.prog)
	in := make(chan int64)
	defer close(in)
	out := intcode.Run(prog, in, debug)
	in <- int64(x)
	in <- int64(y)
	return state(<-out)
}

type beamReadings struct {
	width, height int
	m             map[geom.Pt2]state
}

func mapBeamReadings(prog []int64, width, height int) beamReadings {
	d := newDrone(prog)
	m := beamReadings{width, height, make(map[geom.Pt2]state)}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			m.m[geom.Pt2{x, y}] = d.test(x, y, false)
		}
	}
	return m
}

func printBeamReadings(m beamReadings) {
	for x := 0; x < m.width; x++ {
		for y := 0; y < m.height; y++ {
			switch m.m[geom.Pt2{x, y}] {
			case pulled:
				fmt.Print("#")
			case stationary:
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func countBeamReadings(m beamReadings) int {
	affected := 0
	for _, v := range m.m {
		if v == pulled {
			affected++
		}
	}
	return affected
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	m := mapBeamReadings(data, 50, 50)
	printBeamReadings(m)
	fmt.Println(countBeamReadings(m))
}
