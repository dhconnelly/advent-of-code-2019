package geom

import "testing"

func TestAdd(t *testing.T) {
	for _, tc := range []struct {
		pt1, pt2, want Pt2
	}{
		{Zero2, Zero2, Zero2},
		{Pt2{1, 2}, Pt2{-2, -1}, Pt2{-1, 1}},
		{Pt2{1, 2}, Pt2{2, 1}, Pt2{3, 3}},
	} {
		got := tc.pt1.Add(tc.pt2)
		if got != tc.want {
			t.Errorf("Add(%d, %d) = %d, want %d", tc.pt1, tc.pt2, got, tc.want)
		}
	}
}

func TestManhattanDist(t *testing.T) {
	for _, tc := range []struct {
		pt1, pt2 Pt2
		dist     int
	}{
		{Zero2, Zero2, 0},
		{Zero2, Pt2{1, 1}, 2},
		{Pt2{-2, -3}, Zero2, 5},
		{Pt2{-3, 5}, Pt2{4, -7}, 19},
	} {
		dist := tc.pt1.ManhattanDist(tc.pt2)
		if dist != tc.dist {
			t.Errorf("ManhattanDist(%v, %v) = %d, want %d", tc.pt1, tc.pt2, dist, tc.dist)
		}
	}
}

func TestManhattanNorm(t *testing.T) {
	for _, tc := range []struct {
		pt   Pt2
		dist int
	}{
		{Zero2, 0},
		{Pt2{1, 1}, 2},
		{Pt2{-2, -3}, 5},
		{Pt2{4, -7}, 11},
	} {
		dist := tc.pt.ManhattanNorm()
		if dist != tc.dist {
			t.Errorf("ManhattanNorm(%v) = %d, want %d", tc.pt, dist, tc.dist)
		}
	}
}
