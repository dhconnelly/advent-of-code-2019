package main

import (
	"fmt"
	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"github.com/dhconnelly/advent-of-code-2019/ints"
	"log"
	"math"
	"os"
)

type color int64

const (
	BLACK color = iota
	WHITE
)

type direction int64

const (
	TURN_LEFT direction = iota
	TURN_RIGHT
)

type orientation rune

const (
	LEFT orientation = iota
	UP
	RIGHT
	DOWN
)

func turn(cur orientation, dir direction) orientation {
	switch dir {
	case TURN_LEFT:
		if cur == LEFT {
			return DOWN
		}
		return cur - 1
	case TURN_RIGHT:
		return (cur + 1) % 4
	}
	log.Fatalf("bad dir: %d", dir)
	return -1
}

var diffs = map[orientation]geom.Pt2{
	LEFT:  geom.Pt2{-1, 0},
	UP:    geom.Pt2{0, 1},
	RIGHT: geom.Pt2{1, 0},
	DOWN:  geom.Pt2{0, -1},
}

func move(cur geom.Pt2, o orientation) geom.Pt2 {
	return cur.Add(diffs[o])
}

type grid map[geom.Pt2]color

func run(data []int64, initial color) grid {
	in := make(chan int64)
	out := intcode.RunProgram(data, in)
	g := grid(make(map[geom.Pt2]color))
	p := geom.Zero2
	g[p] = initial
	o := UP
loop:
	for {
		select {
		case c, ok := <-out:
			if !ok {
				break loop
			}
			g[p] = color(c)
			dir := direction(<-out)
			o = turn(o, dir)
			p = move(p, o)
		case in <- int64(g[p]):
		}
	}
	return g
}

func printGrid(g grid) {
	minX, minY := math.MaxInt64, math.MaxInt64
	maxX, maxY := math.MinInt64, math.MinInt64
	for p, _ := range g {
		minX, maxX = ints.Min(minX, p.X), ints.Max(maxX, p.X)
		minY, maxY = ints.Min(minY, p.Y), ints.Max(maxY, p.Y)
	}
	for row := maxY; row >= minY; row-- {
		for col := minX; col <= maxX; col++ {
			p := geom.Pt2{col, row}
			switch g[p] {
			case BLACK:
				fmt.Print("  ")
			case WHITE:
				fmt.Print("XX")
			}
		}
		fmt.Println()
	}
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	g := run(data, BLACK)
	fmt.Println(len(g))
	g = run(data, WHITE)
	printGrid(g)
}
