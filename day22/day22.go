package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type cut struct {
	d   int
	mod int
}

func (c cut) apply(n int) int {
	n -= c.d
	for n < 0 {
		n += c.mod
	}
	return n % c.mod
}

type deal struct {
	d   int
	mod int
}

func (d deal) apply(n int) int {
	return (n * d.d) % d.mod
}

type redeal struct {
	mod int
}

func (r redeal) apply(n int) int {
	return r.mod - n - 1
}

type tf interface {
	apply(n int) int
}

func ReadTransformations(r io.Reader, mod int) []tf {
	scan := bufio.NewScanner(r)
	var tfs []tf
	for scan.Scan() {
		line := scan.Text()
		switch {
		case strings.Index(line, "deal with") == 0:
			d := deal{mod: mod}
			fmt.Sscanf(line, "deal with increment %d", &d.d)
			tfs = append(tfs, d)
		case strings.Index(line, "cut") == 0:
			c := cut{mod: mod}
			fmt.Sscanf(line, "cut %d", &c.d)
			tfs = append(tfs, c)
		case strings.Index(line, "deal into") == 0:
			r := redeal{mod: mod}
			tfs = append(tfs, r)
		default:
			log.Fatal("bad line:", line)
		}
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
	return tfs
}

func apply(n int, tfs []tf) int {
	for _, t := range tfs {
		n = t.apply(n)
	}
	return n
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	tfs := ReadTransformations(f, 10007)
	fmt.Println(tfs)
	fmt.Println(apply(2019, tfs))
}
