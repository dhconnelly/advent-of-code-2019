package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dhconnelly/advent-of-code-2019/intcode"
)

type packet struct {
	dest int64
	x, y int64
}

type machine struct {
	addr int64
	q    []packet
	in   chan<- int64
	out  <-chan int64
	snd  chan<- packet
	rcv  chan packet
	last int64
}

func (m *machine) run() {
	for {
		if len(m.q) > 0 {
			first := m.q[0]
			select {
			case p := <-m.rcv:
				if first.x == -1 {
					m.q = m.q[1:]
				}
				m.q = append(m.q, p)
			case dest := <-m.out:
				x := <-m.out
				y := <-m.out
				m.snd <- packet{dest, x, y}
			case m.in <- first.x:
				m.q = m.q[1:]
				if first.x != -1 {
					m.in <- first.y
					if len(m.q) == 0 {
						m.q = append(m.q, packet{x: -1})
					}
				}
			}
		} else {
			select {
			case p := <-m.rcv:
				m.q = append(m.q, p)
			case dest := <-m.out:
				x := <-m.out
				y := <-m.out
				m.snd <- packet{dest, x, y}
			}
		}
	}
}

func newMachine(
	addr int64,
	prog []int64,
	snd chan<- packet,
) *machine {
	in := make(chan int64, 1)
	in <- addr
	rcv := make(chan packet)
	out := intcode.RunProgram(prog, in)
	m := &machine{
		addr: addr,
		q:    []packet{{x: -1}},
		rcv:  rcv,
		snd:  snd,
		in:   in,
		out:  out,
	}
	go m.run()
	return m
}

func machines(
	n int,
	prog []int64,
	out chan<- packet,
) map[int64]*machine {
	ms := make(map[int64]*machine)
	for i := int64(0); i < int64(n); i++ {
		ms[i] = newMachine(i, prog, out)
	}
	return ms
}

func network(n int, prog []int64) int64 {
	out := make(chan packet)
	ms := machines(n, prog, out)
	for p := range out {
		fmt.Println(p)
		if p.dest == 255 {
			return p.y
		}
		ms[p.dest].rcv <- p
	}
	log.Fatal("network failed")
	return 0
}

func main() {
	prog, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(network(50, prog))
}
