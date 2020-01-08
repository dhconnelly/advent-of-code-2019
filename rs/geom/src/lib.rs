struct Point2 {
    x: i64,
    y: i64,
}

impl Point2 {
    fn manhattan_dist(&self, q: &Point2) -> i64 {
        (self.x - q.x).abs() + (self.y - q.y).abs()
    }

    fn manhattan_norm(&self) -> i64 {
        self.manhattan_dist(&Point2 { x: 0, y: 0 })
    }
}

#[test]
fn manhattan_dist() {
    let p = Point2 { x: 4, y: -5 };
    let q = Point2 { x: -10, y: 7 };
    assert_eq!(p.manhattan_dist(&q), 26);
}

#[test]
fn manhattan_norm() {
    let p = Point2 { x: 4, y: -5 };
    assert_eq!(p.manhattan_norm(), 9);
}
