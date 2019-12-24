package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type bitset int64

func (b bitset) get(i int) bool {
	return (b & (1 << i)) > 0
}

func (b *bitset) set(i int, value bool) {
	if value {
		(*b) |= 1 << i
	} else {
		(*b) &= ^(1 << i)
	}
}

type layout struct {
	bits   bitset
	width  int
	height int
}

func (l *layout) alive(row, col int) bool {
	return l.bits.get(row*l.width + col)
}

func (l *layout) adjacent(row, col int) int {
	adj := 0
	if row > 0 && l.alive(row-1, col) {
		adj++
	}
	if row < l.height-1 && l.alive(row+1, col) {
		adj++
	}
	if col > 0 && l.alive(row, col-1) {
		adj++
	}
	if col < l.width-1 && l.alive(row, col+1) {
		adj++
	}
	return adj
}

func (l *layout) next() {
	next := l.bits
	for row := 0; row < l.height; row++ {
		for col := 0; col < l.width; col++ {
			adj := l.adjacent(row, col)
			n := row*l.width + col
			if l.bits.get(n) && adj != 1 {
				next.set(n, false)
			} else if !l.bits.get(n) && (adj == 1 || adj == 2) {
				next.set(n, true)
			}
		}
	}
	l.bits = next
}

func (l layout) String() string {
	n := l.width * l.height
	b := make([]byte, n)
	for j := 0; j < n; j++ {
		if l.bits.get(j) {
			b[j] = '#'
		} else {
			b[j] = '.'
		}
	}
	return string(b)
}

func readLayout(r io.Reader) layout {
	b := make([]byte, 1)
	var l layout
	i := 0
outer:
	for {
	inner:
		for {
			_, err := r.Read(b)
			if err == io.EOF {
				break outer
			}
			if err != nil {
				log.Fatal(err)
			}
			switch b[0] {
			case '\n':
				break inner
			case '#':
				l.bits.set(i, true)
				i++
			case '.':
				l.bits.set(i, false)
				i++
			default:
				log.Fatalf("bad char in layout: %c", b[0])
			}
		}
		l.height++
	}
	l.width = i / l.height
	return l
}

func findRepeat(l layout) bitset {
	m := map[bitset]bool{l.bits: true}
	for {
		l.next()
		if _, ok := m[l.bits]; ok {
			return l.bits
		}
		m[l.bits] = true
	}
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	l := readLayout(f)
	fmt.Println(findRepeat(l))
}
