package geom

type Rect struct {
	Lo, Hi Pt2 // inclusive
}

func (r Rect) Contains(p Pt2) bool {
	return p.X >= r.Lo.X && p.Y >= r.Lo.Y && p.X <= r.Hi.X && p.Y <= r.Hi.Y
}
