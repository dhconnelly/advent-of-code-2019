package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/dhconnelly/advent-of-code-2019/intcode"
)

const idleReceiveCount = 100

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

type idleMap struct {
	sync.RWMutex
	count int64
	idles map[int64]int64
}

func NewIdleMap(count int64) *idleMap {
	m := make(map[int64]int64)
	return &idleMap{idles: m, count: count}
}

func (im *idleMap) incr(addr int64) {
	im.Lock()
	defer im.Unlock()
	im.idles[addr]++
}

func (im *idleMap) reset(addr int64) {
	im.Lock()
	defer im.Unlock()
	im.idles[addr] = 0
}

func (im *idleMap) allIdle() bool {
	im.RLock()
	defer im.RUnlock()
	for addr := int64(0); addr < im.count; addr++ {
		if im.idles[addr] < idleReceiveCount {
			return false
		}
	}
	return true
}

func NewMachine(
	prog []int64,
	addr int64,
	idles *idleMap,
) machine {
	out, in := make(chan packet), make(chan packet)
	min := make(chan int64)
	mout := intcode.RunProgram(prog, min)
	min <- int64(addr)
	var q []packet

	go func() {
		for {
			select {
			case to, ok := <-mout:
				if idles != nil {
					idles.reset(addr)
				}
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
				if idles != nil {
					idles.reset(addr)
				}

			case min <- head(q):
				if head(q) == -1 {
					if idles != nil {
						idles.incr(addr)
					}
				} else {
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

type nat struct {
	last packet
	m    machine
}

func NewNat(addr int64, idles *idleMap) nat {
	in := make(chan packet)
	out := make(chan packet)
	n := nat{m: machine{addr: addr, in: in, out: out}}
	go func() {
		for {
			select {
			case p := <-in:
				fmt.Println(p)
				n.last = p
			default:
				if idles.allIdle() {
					fmt.Println("all idle, sending", n.last)
					idles.reset(0)
					out <- packet{from: addr, to: 0, x: n.last.x, y: n.last.y}
				}
			}
		}
	}()
	return n
}

func networkSwitch(ms map[int64]machine, n *nat) packet {
	out := make(chan packet)
	closed := make(chan int64)
	closedCount := 0
	for addr, m := range ms {
		go listen(addr, m, out, closed)
	}
	var lastNatPacket packet

loop:
	for {
		select {
		case p, ok := <-out:
			if !ok {
				break loop
			}
			if n == nil && p.to == 255 {
				return p
			}
			if n != nil && p.from == 255 && p.to == 0 {
				if lastNatPacket.y == p.y {
					return lastNatPacket
				}
				lastNatPacket = p
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

func machines(
	data []int64,
	count int64,
	idles *idleMap,
) map[int64]machine {
	ms := make(map[int64]machine)
	for addr := int64(0); addr < count; addr++ {
		ms[addr] = NewMachine(data, addr, idles)
	}
	return ms
}

func main() {
	data, err := intcode.ReadProgram(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	ms := machines(data, 50, nil)
	fmt.Println(networkSwitch(ms, nil).y)

	idles := NewIdleMap(50)
	n := NewNat(255, idles)
	ms = machines(data, 50, idles)
	ms[255] = n.m
	fmt.Println(networkSwitch(ms, &n).y)
}
