package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
	"github.com/dhconnelly/advent-of-code-2019/ints"
)

type elem rune

const SPACE elem = '.'

type grid struct {
	height, width int
	g             map[geom.Pt2]elem
}

func (g grid) neighbors(p geom.Pt2) []geom.Pt2 {
	var nbrs []geom.Pt2
	for _, nbr := range p.ManhattanNeighbors() {
		if c, ok := g.g[nbr]; ok && c != '.' {
			nbrs = append(nbrs, nbr)
		}
	}
	return nbrs
}

func readGridFrom(out <-chan int64) (grid, bool) {
	g := grid{g: make(map[geom.Pt2]elem)}
	var width int
	for ch := range out {
		if ch == '\n' {
			if width > 0 {
				g.height++
				g.width = width
				width = 0
				continue
			} else {
				return g, true
			}
		}
		g.g[geom.Pt2{width, g.height}] = elem(ch)
		width++
	}
	return grid{}, false
}

func readGrid(data []int64) grid {
	data = ints.Copied64(data)
	in := make(chan int64)
	out := intcode.RunProgram(data, in)
	g, ok := readGridFrom(out)
	if !ok {
		log.Fatal("can't read grid")
	}
	return g
}

func readGraph(g grid) map[geom.Pt2][]geom.Pt2 {
	m := make(map[geom.Pt2][]geom.Pt2)
	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			p := geom.Pt2{j, i}
			if c := g.g[p]; c == '.' {
				continue
			}
			var edges []geom.Pt2
			for _, nbr := range g.neighbors(p) {
				edges = append(edges, nbr)
			}
			m[p] = edges
		}
	}
	return m
}

func intersections(m map[geom.Pt2][]geom.Pt2) []geom.Pt2 {
	var ps []geom.Pt2
	for p, edges := range m {
		if len(edges) > 2 {
			ps = append(ps, p)
		}
	}
	return ps
}

func alignmentSum(g grid) int {
	m := readGraph(g)
	ps := intersections(m)
	sum := 0
	for _, p := range ps {
		sum += p.X * p.Y
	}
	return sum
}

func writeLine(ch chan<- int64, line string) {
	for _, c := range line {
		ch <- int64(c)
	}
	ch <- int64('\n')
}

func readLine(ch <-chan int64) string {
	var s []rune
	for {
		c := <-ch
		if c == '\n' {
			return string(s)
		}
		s = append(s, rune(c))
	}
}

func computeDust(data []int64, prog [4]string) int64 {
	data = ints.Copied64(data)
	data[0] = 2
	in := make(chan int64)
	out := intcode.RunProgram(data, in)
	readGridFrom(out)
	for _, line := range prog {
		readLine(out)
		writeLine(in, line)
	}
	readLine(out)
	writeLine(in, "n")
	var answer int64
	for c := range out {
		answer = c
	}
	return answer
}

var prog = [4]string{
	"A,B,A,C,B,A,C,B,A,C",
	"L,12,L,12,L,6,L,6",
	"R,8,R,4,L,12",
	"L,12,L,6,R,12,R,8",
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	g := readGrid(data)
	fmt.Println(alignmentSum(g))
	fmt.Println(computeDust(data, prog))
}
