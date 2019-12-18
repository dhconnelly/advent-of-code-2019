package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/geom"
)

const (
	door = '@'
	wall = '#'
	open = '.'
)

func isDoor(c rune) bool {
	return 'A' <= c && c <= 'Z'
}

func isKey(c rune) bool {
	return 'a' <= c && c <= 'z'
}

type maze struct {
	width  int
	height int
	grid   map[geom.Pt2]rune
}

func (m maze) find(c rune) geom.Pt2 {
	for k, v := range m.grid {
		if v == c {
			return k
		}
	}
	log.Fatalf("not found in maze: %c", c)
	return geom.Zero2
}

func (m maze) adjacent(p geom.Pt2) []geom.Pt2 {
	var adj []geom.Pt2
	for _, q := range p.ManhattanNeighbors() {
		if c, ok := m.grid[q]; ok && c != wall {
			adj = append(adj, q)
		}
	}
	return adj
}

func readMaze(r io.Reader) maze {
	scan := bufio.NewScanner(r)
	m := maze{grid: make(map[geom.Pt2]rune)}
	for scan.Scan() {
		m.width = 0
		for _, c := range scan.Text() {
			m.grid[geom.Pt2{m.width, m.height}] = c
			m.width++
		}
		m.height++
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
	return m
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m := readMaze(f)
	p := m.find(door)
	fmt.Println(m)
	fmt.Println(p)
	fmt.Println(m.adjacent(m.adjacent(p)[0]))
}
