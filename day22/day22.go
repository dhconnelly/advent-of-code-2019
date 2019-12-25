package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strings"
)

var one = big.NewInt(1)
var minusOne = big.NewInt(-1)
var zero = big.NewInt(0)

type applier interface {
	apply(n *big.Int, mod *big.Int)
}

type change struct {
	scale *big.Int
	shift *big.Int
}

func (chg change) mod(mod int64) {
	chg.scale.Mod(chg.scale, big.NewInt(mod))
	chg.shift.Mod(chg.shift, big.NewInt(mod))
}

func (chg change) apply(n *big.Int, mod *big.Int) {
	n.Mul(n, chg.scale)
	n.Add(n, chg.shift)
	n.Mod(n, mod)
}

func (chg change) invert(mod int64) invertedChange {
	chg.mod(mod)
	m := big.NewInt(mod)
	scale := big.NewInt(0)
	if chg.scale.Cmp(zero) != 0 {
		scale = chg.scale.ModInverse(chg.scale, m)
	}
	shift := chg.shift.Mul(chg.shift, minusOne)
	shift.Mul(shift, scale)
	shift.Mod(shift, m)
	return invertedChange{scale, shift}
}

type invertedChange struct {
	scale *big.Int
	shift *big.Int
}

func (ichg invertedChange) apply(n *big.Int, mod *big.Int) {
	n.Mul(n, ichg.scale)
	n.Add(n, ichg.shift)
	n.Mod(n, mod)
}

func geoSum(a, r, n, mod *big.Int) {
	var denom big.Int
	denom.Set(one)
	denom.Sub(&denom, r)
	denom.ModInverse(&denom, mod)

	r.Exp(r, n, mod)
	r.Mul(r, minusOne)
	r.Add(r, one)

	r.Mul(r, &denom)
	r.Mul(r, a)

	a.Set(r)
	a.Mod(a, mod)
}

func (ichg invertedChange) pow(n *big.Int, mod *big.Int) {
	var scale big.Int
	scale.Set(ichg.scale)
	geoSum(ichg.shift, &scale, n, mod)
	ichg.scale.Exp(ichg.scale, n, mod)
}

func ReadTransformations(r io.Reader) change {
	scan := bufio.NewScanner(r)
	chg := change{big.NewInt(1), big.NewInt(0)}
	var x big.Int
	for scan.Scan() {
		line := scan.Text()
		switch {
		case strings.Index(line, "deal with") == 0:
			var d int64
			fmt.Sscanf(line, "deal with increment %d", &d)
			x.SetInt64(d)
			chg.scale.Mul(chg.scale, &x)
			chg.shift.Mul(chg.shift, &x)
		case strings.Index(line, "cut") == 0:
			var d int64
			fmt.Sscanf(line, "cut %d", &d)
			x.SetInt64(d)
			chg.shift.Sub(chg.shift, &x)
		case strings.Index(line, "deal into") == 0:
			chg.scale.Mul(chg.scale, minusOne)
			chg.shift.Mul(chg.shift, minusOne)
			chg.shift.Sub(chg.shift, one)
		default:
			log.Fatal("bad line:", line)
		}
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
	return chg
}

func apply(n int64, mod int64, a applier, times int) int64 {
	x := big.NewInt(n)
	m := big.NewInt(mod)
	for ; times > 0; times-- {
		a.apply(x, m)
		x.Mod(x, m)
	}
	return x.Int64()
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	chg := ReadTransformations(f)
	fmt.Println(apply(2019, 10007, chg, 1))

	var mod int64 = 119315717514047
	const n = 2020
	const times = 101741582076661
	i := chg.invert(mod)
	m := big.NewInt(mod)
	x := big.NewInt(n)
	i.pow(big.NewInt(times), m)
	i.apply(x, m)
	x.Mod(x, m)
	fmt.Println(x)
}
