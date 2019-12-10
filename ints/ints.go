package ints

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Gcd(m, n int) int {
	if n == 0 {
		return m
	}
	return Gcd(n, m%n)
}
