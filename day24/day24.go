package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type layout struct {
	bits   int64
	width  int
	height int
}

func (l layout) alive(x, y int) bool {
	i := y*l.width + x
	return (l.bits & (1 << i)) > 0
}

func (l layout) String() string {
	i, n := l.bits, l.width*l.height
	b := make([]byte, n)
	for j := 0; j < n; j++ {
		if (i>>j)&1 == 1 {
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
				l.bits |= (1 << i)
				i++
			case '.':
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

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	l := readLayout(f)
	fmt.Println(l)
	fmt.Printf("%#v\n", l)
	for j := 0; j < 5; j++ {
		for i := 0; i < 5; i++ {
			fmt.Println(l.alive(i, j))
		}
	}
}
