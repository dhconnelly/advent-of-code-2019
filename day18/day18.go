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

func (m maze) at(p geom.Pt2) rune {
	return m.grid[p]
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

type explorer struct {
	m     maze
	keys  map[rune]bool
	path  []rune
	steps int
}

type node struct {
	p geom.Pt2
	c rune
	d int
}

func (nd node) String() string {
	return fmt.Sprintf("(%v, %c, dist=%d)", nd.p, nd.c, nd.d)
}

func (e explorer) reachableKeys(from geom.Pt2) []node {
	// breadth-first-search from current point
	var keys []node
	var nd node
	visited := make(map[geom.Pt2]bool)
	q := []node{{from, e.m.at(from), 0}}
	for len(q) > 0 {
		nd, q = q[0], q[1:]
		visited[nd.p] = true

		// if we're at a key or a door, go no further
		if c := e.m.at(nd.p); isKey(c) {
			keys = append(keys, nd)
			continue
		} else if isDoor(c) {
			continue
		}

		// otherwise keep going
		for _, nbr := range e.m.adjacent(nd.p) {
			if !visited[nbr] {
				q = append(q, node{nbr, e.m.at(nbr), nd.d + 1})
			}
		}
	}
	return keys
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m := readMaze(f)
	p := m.find(door)
	e := explorer{m: m, keys: make(map[rune]bool)}
	fmt.Println(e.reachableKeys(p))
	fmt.Println(e.reachableKeys(geom.Pt2{4, 3}))
	fmt.Println(e.reachableKeys(geom.Pt2{7, 7}))
}
