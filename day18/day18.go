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

func doorFor(key rune) rune {
	return key - 32
}

func keyFor(door rune) rune {
	return door + 32
}

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

func (m maze) clone() maze {
	m2 := maze{m.width, m.height, make(map[geom.Pt2]rune)}
	for k, v := range m.grid {
		m2.grid[k] = v
	}
	return m2
}

func (m maze) keys() []rune {
	var keys []rune
	for _, v := range m.grid {
		if isKey(v) {
			keys = append(keys, v)
		}
	}
	return keys
}

func (m *maze) clear(p geom.Pt2) {
	m.grid[p] = open
}

func (m maze) at(p geom.Pt2) rune {
	return m.grid[p]
}

func (m maze) find(c rune) (geom.Pt2, bool) {
	for k, v := range m.grid {
		if v == c {
			return k, true
		}
	}
	return geom.Zero2, false
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

type keyPath struct {
	path  []rune
	steps int
}

func findKeyPaths(m maze, from geom.Pt2) []keyPath {
	var paths []keyPath
	ch := make(chan keyPath)
	e := explorer{m: m, keys: make(map[rune]bool), out: ch}
	go e.findKeys(from)
	for path := range ch {
		fmt.Println("path:", path)
		paths = append(paths, path)
	}
	return paths
}

type explorer struct {
	m    maze
	keys map[rune]bool
	path keyPath
	out  chan<- keyPath
}

func (e explorer) clone() explorer {
	e2 := explorer{m: e.m.clone(), out: e.out}
	e2.keys = make(map[rune]bool)
	for k, v := range e.keys {
		e2.keys[k] = v
	}
	e2.path.steps = e.path.steps
	e2.path.path = make([]rune, len(e.path.path))
	for i, c := range e.path.path {
		e2.path.path[i] = c
	}
	return e2
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

func (e *explorer) findKeys(from geom.Pt2) {
	//fmt.Println("findKeys:", from)
	//fmt.Println("have:", e.keys)
	//fmt.Println("remaining:", e.m.keys())

	// if no more keys remaining, send path on channel
	if len(e.m.keys()) == 0 {
		e.out <- e.path
		return
	}

	// find current reachable keys
	reachable := e.reachableKeys(from)

	// if none, nothing to do
	if len(reachable) == 0 {
		//fmt.Println("none reachable! dying")
		return
	}

	// if more than one, fork and choose one per clone
	if len(reachable) > 1 {
		for _, nd := range reachable[1:] {
			//fmt.Println("forking")
			clone := e.clone()
			clone.takeKey(nd)
			go clone.findKeys(nd.p)
		}
	}

	// take the first in this clone
	nd := reachable[0]
	e.takeKey(nd)

	// continue from the key location
	e.findKeys(nd.p)
}

func (e *explorer) takeKey(nd node) {
	e.path.path = append(e.path.path, nd.c)
	e.path.steps += nd.d
	e.keys[nd.c] = true
	e.m.clear(nd.p)
	door, ok := e.m.find(doorFor(nd.c))
	if ok {
		e.m.clear(door)
	}
}

func values(m map[geom.Pt2]node) []node {
	var vs []node
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

func reachable(m maze, from geom.Pt2) []node {
	var nd node
	nds := make(map[geom.Pt2]node)
	visited := make(map[geom.Pt2]bool)
	q := []node{{from, m.at(from), 0}}
	for len(q) > 0 {
		nd, q = q[0], q[1:]
		visited[nd.p] = true
		for _, nbr := range m.adjacent(nd.p) {
			if visited[nbr] {
				continue
			}
			next := node{nbr, m.at(nbr), nd.d + 1}
			if c := m.at(nbr); isKey(c) || isDoor(c) {
				if prev, ok := nds[nbr]; !ok || next.d < prev.d {
					nds[nbr] = next
				}
				continue
			}
			q = append(q, next)
		}
	}
	return values(nds)
}

func reachability(m maze, from geom.Pt2) map[rune][]rune {
	g := make(map[rune][]rune)
	q := []geom.Pt2{from}
	for len(q) > 0 {
		from, q = q[0], q[1:]
		c := m.at(from)
		if _, ok := g[c]; ok {
			continue
		}
		for _, nbr := range reachable(m, from) {
			g[c] = append(g[c], nbr.c)
			q = append(q, nbr.p)
		}
	}
	return g
}

func printReachability(m map[rune][]rune) {
	for k, v := range m {
		fmt.Printf("[%c]: %s\n", k, string(v))
	}
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m := readMaze(f)
	p, _ := m.find(door)
	g := reachability(m, p)
	printReachability(g)
}
