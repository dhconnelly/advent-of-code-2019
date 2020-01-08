enum Direction {
    Up,
    Down,
    Left,
    Right,
}

use std::f64::consts::PI;

#[derive(Clone, Copy, Debug, PartialEq)]
struct Point2 {
    x: i32,
    y: i32,
}

const ZERO2: Point2 = Point2 { x: 0, y: 0 };

impl Point2 {
    fn manhattan_dist(&self, q: &Point2) -> i32 {
        (self.x - q.x).abs() + (self.y - q.y).abs()
    }

    fn manhattan_norm(&self) -> i32 {
        self.manhattan_dist(&Point2 { x: 0, y: 0 })
    }

    fn go(&self, dir: Direction) -> Point2 {
        let mut q = *self;
        match dir {
            Direction::Up => q.y += 1,
            Direction::Down => q.y -= 1,
            Direction::Left => q.x -= 1,
            Direction::Right => q.x += 1,
        }
        q
    }

    fn manhattan_neighbors(&self) -> [Point2; 4] {
        [
            self.go(Direction::Up),
            self.go(Direction::Down),
            self.go(Direction::Left),
            self.go(Direction::Right),
        ]
    }

    fn angle_from(&self, q: &Point2) -> f64 {
        let x = (self.x - q.x) as f64;
        let y = (self.y - q.y) as f64;
        y.atan2(x)
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

#[test]
fn go() {
    let p = Point2 { x: 4, y: -5 };
    assert_eq!(p.go(Direction::Up), Point2 { x: 4, y: -4 });
    assert_eq!(p.go(Direction::Down), Point2 { x: 4, y: -6 });
    assert_eq!(p.go(Direction::Left), Point2 { x: 3, y: -5 });
    assert_eq!(p.go(Direction::Right), Point2 { x: 5, y: -5 });
}

#[test]
fn manhattan_neighbors() {
    let p = Point2 { x: 4, y: -5 };
    let nbrs = p.manhattan_neighbors();
    for q in nbrs.iter() {
        assert_eq!(p.manhattan_dist(q), 1);
    }
}

#[test]
fn angle_from() {
    let p = Point2 { x: 1, y: 1 };
    assert_eq!(p.angle_from(&ZERO2), PI / 4.0);
    assert_eq!(ZERO2.angle_from(&p), -3.0 * PI / 4.0);
}
