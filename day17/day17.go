package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/geom"
	"github.com/dhconnelly/advent-of-code-2019/intcode"
)

type elem rune

const (
	SPACE elem = '.'
	SCAFF elem = '#'
	DROID elem = '^'
)

type grid struct {
	height, width int
	g             map[geom.Pt2]elem
}

func (g grid) neighbors(p geom.Pt2) []geom.Pt2 {
	var nbrs []geom.Pt2
	for _, nbr := range p.ManhattanNeighbors() {
		if c, ok := g.g[nbr]; ok && c != SPACE {
			nbrs = append(nbrs, nbr)
		}
	}
	return nbrs
}

func readGrid(data []int64) grid {
	in := make(chan int64)
	out := intcode.RunProgram(data, in)
	g := grid{g: make(map[geom.Pt2]elem)}
	var width int
	for ch := range out {
		if ch == '\n' {
			if width > 0 {
				g.height++
				g.width = width
				width = 0
			}
			continue
		}
		g.g[geom.Pt2{width, g.height}] = elem(ch)
		width++
	}
	return g
}

func readGraph(g grid) map[geom.Pt2][]geom.Pt2 {
	m := make(map[geom.Pt2][]geom.Pt2)
	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			p := geom.Pt2{j, i}
			if c := g.g[p]; c == SPACE {
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

func printGrid(g grid) {
	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			p := geom.Pt2{j, i}
			c := g.g[p]
			if c != SPACE && len(g.neighbors(p)) > 2 {
				fmt.Printf("O")
			} else {
				fmt.Printf("%c", c)
			}
		}
		fmt.Println()
	}
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
	fmt.Println(ps)
	sum := 0
	for _, p := range ps {
		sum += p.X * p.Y
	}
	return sum
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	g := readGrid(data)
	printGrid(g)
	fmt.Println(alignmentSum(g))
}
