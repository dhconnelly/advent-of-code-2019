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

func NewMachine(prog []int64, addr int64) machine {
	out, in := make(chan packet), make(chan packet)
	min := make(chan int64)
	mout := intcode.RunProgram(prog, min)
	min <- int64(addr)
	var q []packet
	var head int64 = -1
	go func() {
		for {
			select {
			case to, ok := <-mout:
				log.Printf("machine %d sending to %d", addr, to)
				if !ok {
					close(in)
					close(out)
					close(min)
					return
				}
				x, y := <-mout, <-mout
				out <- packet{addr, to, x, y}
			case p := <-in:
				log.Printf("queuing %v for machine %d", p, addr)
				q = append(q, p)
				head = p.x
			case min <- head:
				log.Printf("machine %d receiving %d", addr, head)
				if head != -1 {
					p := q[0]
					q = q[1:]
					min <- p.y
					head = -1
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
	for {
		p, ok := <-m.out
		if !ok {
			closed <- addr
			return
		}
		out <- p
	}
}

func networkSwitch(ms map[int64]machine) packet {
	out := make(chan packet)
	closed := make(chan int64)
	closedCount := 0
	for addr, m := range ms {
		go listen(addr, m, out, closed)
	}
loop:
	for {
		select {
		case addr := <-closed:
			log.Println("machine halted:", addr)
			closedCount++
			if closedCount == len(ms) {
				close(out)
				close(closed)
			}
		case p, ok := <-out:
			log.Println("packet:", p)
			if !ok {
				break loop
			}
			if p.to == 255 {
				return p
			}
			ms[p.to].in <- p
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
	fmt.Println(networkSwitch(ms))
}
