package main

import (
	"bufio"
	"fmt"
	"github.com/dhconnelly/advent-of-code-2019/ints"
	"log"
	"os"
)

type state struct {
	px, py, pz, vx, vy, vz [4]int64
}

func readState(path string) state {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	var s state
	i := 0
	for scan.Scan() {
		line := scan.Text()
		_, err = fmt.Sscanf(
			line, "<x=%d, y=%d, z=%d>",
			&s.px[i], &s.py[i], &s.pz[i],
		)
		if err != nil {
			log.Fatalf("bad point: %s", err)
		}
		i++
	}
	if i > 4 {
		log.Fatalf("bad number of moons: %d", i)
	}
	if err = scan.Err(); err != nil {
		log.Fatal(err)
	}
	return s
}

func applyGravity(px, vx *[4]int64) {
	for i := 0; i < len(px)-1; i++ {
		for j := i + 1; j < len(px); j++ {
			if px[i] < px[j] {
				vx[i] += 1
				vx[j] -= 1
			} else if px[i] > px[j] {
				vx[i] -= 1
				vx[j] += 1
			}
		}
	}
}

func applyVelocity(px, vx *[4]int64) {
	for i := range px {
		px[i] += vx[i]
	}
}

func step(px, vx *[4]int64) {
	applyGravity(px, vx)
	applyVelocity(px, vx)
}

func simulate(s state, n int) state {
	for i := 0; i < n; i++ {
		step(&s.px, &s.vx)
		step(&s.py, &s.vy)
		step(&s.pz, &s.vz)
	}
	return s
}

func energy(s state) int64 {
	e := int64(0)
	for i := 0; i < len(s.px); i++ {
		pe := ints.Abs64(s.px[i])
		pe += ints.Abs64(s.py[i])
		pe += ints.Abs64(s.pz[i])
		ke := ints.Abs64(s.vx[i])
		ke += ints.Abs64(s.vy[i])
		ke += ints.Abs64(s.vz[i])
		e += pe * ke
	}
	return e
}

func lcm(a, b, c int64) int64 {
	lcm := a * (b / ints.Gcd64(a, b))
	return c * (lcm / ints.Gcd64(lcm, c))
}

func findLoopCoord(px, vx [4]int64, ch chan<- int64) {
	pi, vi := px, vx
	for i := int64(1); ; i++ {
		step(&px, &vx)
		if px == pi && vx == vi {
			ch <- i
			return
		}
	}
}

func findLoop(s state) int64 {
	ch := make(chan int64)
	defer close(ch)
	go findLoopCoord(s.px, s.vx, ch)
	go findLoopCoord(s.py, s.vy, ch)
	go findLoopCoord(s.pz, s.vz, ch)
	return lcm(<-ch, <-ch, <-ch)
}

func main() {
	s := readState(os.Args[1])
	fmt.Println(energy(simulate(s, 1000)))
	fmt.Println(findLoop(s))
}
