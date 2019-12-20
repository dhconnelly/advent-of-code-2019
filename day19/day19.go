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

func (d drone) test(x, y int, debug bool) state {
	prog := ints.Copied64(d.prog)
	in := make(chan int64)
	defer close(in)
	out := intcode.Run(prog, in, debug)
	in <- int64(x)
	in <- int64(y)
	return state(<-out)
}

type beamReadings struct {
	x, y, width, height int
	m                   map[geom.Pt2]state
}

func mapBeamReadings(prog []int64, x, y, width, height int) beamReadings {
	d := drone{prog}
	m := beamReadings{x, y, width, height, make(map[geom.Pt2]state)}
	for j := x; j < x+width; j++ {
		for i := y; i < y+height; i++ {
			m.m[geom.Pt2{j, i}] = d.test(j, i, false)
		}
	}
	return m
}

func printBeamReadings(m beamReadings) {
	for x := m.x; x < m.x+m.width; x++ {
		for y := m.y; y < m.y+m.height; y++ {
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

func fastTest(x, y int) bool {
	// see asm.txt for how this was determined
	return 14*x*y >= ints.Abs(149*x*x-127*y*y)
}

func testRange(x, y, width, height int) bool {
	return (fastTest(x, y) &&
		fastTest(x+width-1, y) &&
		fastTest(x, y+height-1) &&
		fastTest(x+width-1, y+height-1))
}

func findRange(fromX, fromY, toX, toY, width, height int) (int, int) {
	for x := fromX; x <= toX; x++ {
		for y := fromY; y <= toY; y++ {
			if testRange(x, y, width, height) {
				return x, y
			}
		}
	}
	return 0, 0
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	m := mapBeamReadings(data, 0, 0, 50, 50)
	fmt.Println(countBeamReadings(m))

	// range chosen by dumping the formula in fastTest
	// into desmos.com/calculator
	x, y := findRange(1500, 1600, 2000, 2200, 100, 100)
	fmt.Println(x*10000 + y)
}
