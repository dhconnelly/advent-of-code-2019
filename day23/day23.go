package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/intcode"
)

type packet struct {
	from int64
	to   int64
	x, y int64
}

type machine struct {
	addr int64
	in   chan<- packet
	out  <-chan packet
}

func head(q []packet) int64 {
	if len(q) == 0 {
		return -1
	}
	return q[0].x
}

func NewMachine(prog []int64, addr int64) machine {
	out, in := make(chan packet), make(chan packet)
	min := make(chan int64)
	mout := intcode.RunProgram(prog, min)
	min <- int64(addr)
	var q []packet
	go func() {
		for {
			select {
			case to, ok := <-mout:
				if !ok {
					close(in)
					close(out)
					close(min)
					return
				}
				x := <-mout
				y := <-mout
				p := packet{addr, to, x, y}
				out <- p
			case p := <-in:
				q = append(q, p)
			case min <- head(q):
				if head(q) != -1 {
					p := q[0]
					q = q[1:]
					min <- p.y
				}
			}
		}
	}()
	return machine{addr, in, out}
}

func listen(
	addr int64,
	m machine,
	out chan<- packet,
	closed chan<- int64,
) {
	for p := range m.out {
		out <- p
	}
	closed <- addr
}

func networkSwitch(ms map[int64]machine, nat bool) packet {
	out := make(chan packet)
	closed := make(chan int64)
	closedCount := 0
	for addr, m := range ms {
		go listen(addr, m, out, closed)
	}
loop:
	for {
		select {
		case p, ok := <-out:
			if !ok {
				break loop
			}
			if !nat && p.to == 255 {
				return p
			}
			ms[p.to].in <- p
		case <-closed:
			closedCount++
			if closedCount == len(ms) {
				close(out)
				close(closed)
				break loop
			}
		}
	}
	log.Fatal("stopped without return value")
	return packet{}
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	ms := make(map[int64]machine)
	for addr := int64(0); addr < 50; addr++ {
		ms[addr] = NewMachine(data, addr)
	}
	fmt.Println(networkSwitch(ms, false).y)
}
