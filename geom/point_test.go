package geom

import "testing"

func TestManhattanDist2(t *testing.T) {
	for _, tc := range []struct {
		pt1, pt2 Pt2
		dist     int
	}{
		{Zero2, Zero2, 0},
		{Zero2, Pt2{1, 1}, 2},
		{Pt2{-2, -3}, Zero2, 5},
		{Pt2{-3, 5}, Pt2{4, -7}, 19},
	} {
		dist := ManhattanDist2(tc.pt1, tc.pt2)
		if dist != tc.dist {
			t.Errorf("ManhattanDist2(%d, %d) = %d, want %d", tc.pt1, tc.pt2, dist, tc.dist)
		}
	}
}

func TestManhattanNorm2(t *testing.T) {
	for _, tc := range []struct {
		pt   Pt2
		dist int
	}{
		{Zero2, 0},
		{Pt2{1, 1}, 2},
		{Pt2{-2, -3}, 5},
		{Pt2{4, -7}, 11},
	} {
		dist := ManhattanNorm2(tc.pt)
		if dist != tc.dist {
			t.Errorf("ManhattanNorm2(%d) = %d, want %d", tc.pt, dist, tc.dist)
		}
	}
}
