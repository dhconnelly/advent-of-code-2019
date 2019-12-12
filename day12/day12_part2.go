package main

import (
	"bufio"
	"fmt"
	"github.com/dhconnelly/advent-of-code-2019/ints"
	"log"
	"os"
)

type state struct {
	px [4]int16
	py [4]int16
	pz [4]int16
	vx [4]int16
	vy [4]int16
	vz [4]int16
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

func applyGravity(px, vx *[4]int16) {
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

func applyVelocity(px, vx *[4]int16) {
	for i := range px {
		px[i] += vx[i]
	}
}

func step(px, vx *[4]int16) {
	applyGravity(px, vx)
	applyVelocity(px, vx)
}

type key struct {
	px [4]int16
	vx [4]int16
}

type loop struct {
	n      int64
	px, vx [4]int16
}

func findLoopCoord(px, vx [4]int16) chan loop {
	ch := make(chan loop)
	go func() {
		pi, vi := px, vx
		for i := int64(1); ; i++ {
			step(&px, &vx)
			if px == pi && vx == vi {
				ch <- loop{i, px, vx}
				close(ch)
				return
			}
		}
		close(ch)
	}()
	return ch
}

func lcm3(a, b, c int64) int64 {
	lcm := a * (b / ints.Gcd64(a, b))
	return c * (lcm / ints.Gcd64(lcm, c))
}

func findLoop(s state) int64 {
	ch1 := findLoopCoord(s.px, s.vx)
	ch2 := findLoopCoord(s.py, s.vy)
	ch3 := findLoopCoord(s.pz, s.vz)
	l1, l2, l3 := <-ch1, <-ch2, <-ch3
	return lcm3(l1.n, l2.n, l3.n)
}

func main() {
	s := readState(os.Args[1])
	fmt.Println(findLoop(s))
}
