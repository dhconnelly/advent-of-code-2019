use std::collections::HashMap;
use std::collections::VecDeque;
use std::env;
use std::fs;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
struct Pt2 {
    x: i32,
    y: i32,
}

#[derive(Debug, Clone, Copy, PartialEq)]
enum Position {
    Empty,
    Asteroid,
}

impl Position {
    fn of(ch: char) -> Position {
        match ch {
            '.' => Position::Empty,
            '#' => Position::Asteroid,
            ch => panic!("invalid position: {}", ch),
        }
    }
}

fn read_map(text: &str) -> HashMap<Pt2, Position> {
    let mut map = HashMap::new();
    for (y, xs) in text.lines().enumerate() {
        for (x, ch) in xs.chars().enumerate() {
            map.insert(
                Pt2 {
                    x: x as i32,
                    y: y as i32,
                },
                Position::of(ch),
            );
        }
    }
    map
}

fn asteroids(map: &HashMap<Pt2, Position>) -> Vec<&Pt2> {
    map.iter()
        .filter(|(_, pos)| **pos == Position::Asteroid)
        .map(|(pt, _)| pt)
        .collect()
}

fn dist(from: &Pt2, to: &Pt2) -> f64 {
    ((to.x as f64 - from.x as f64).powi(2) + (to.y as f64 - from.y as f64).powi(2)).sqrt()
}

fn gcd(x: i32, y: i32) -> i32 {
    let x = x.abs();
    let y = y.abs();
    if y == 0 {
        x
    } else {
        gcd(y, x % y)
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
struct Angle {
    dx: i32,
    dy: i32,
}

impl Angle {
    fn of(dx: i32, dy: i32) -> Angle {
        let g = match gcd(dx, dy) {
            0 => 1,
            g => g,
        };
        Angle {
            dx: dx / g,
            dy: dy / g,
        }
    }

    fn between(pt1: &Pt2, pt2: &Pt2) -> Angle {
        Angle::of(pt2.x - pt1.x, pt2.y - pt1.y)
    }

    fn radians(&self) -> f64 {
        let mut a = (self.dx as f64).atan2(self.dy as f64);
        if a <= 0.0 {
            a += 2.0 * std::f64::consts::PI;
        }
        a
    }
}

fn detected_asteroids<'a>(map: &'a HashMap<Pt2, Position>, from: &Pt2) -> Vec<&'a Pt2> {
    let mut detected: HashMap<Angle, &'a Pt2> = HashMap::new();
    for pt in asteroids(map) {
        if pt == from {
            continue;
        }
        let a = Angle::between(pt, from);
        let d = dist(pt, from);
        if let Some(pt1) = detected.get(&a) {
            if d < dist(from, pt1) {
                detected.insert(a, pt);
            }
        } else {
            detected.insert(a, pt);
        }
    }
    detected.values().copied().collect()
}

fn best_position(map: &HashMap<Pt2, Position>) -> &Pt2 {
    asteroids(map)
        .into_iter()
        .map(|pt| (pt, detected_asteroids(map, pt).len()))
        .max_by(|(_, n1), (_, n2)| n1.cmp(n2))
        .unwrap()
        .0
}

fn sort_angular<'a>(pts: &mut Vec<&'a Pt2>, from: &Pt2) {
    pts.sort_by(|pt1, pt2| {
        let a1 = Angle::between(pt1, from).radians();
        let a2 = Angle::between(pt2, from).radians();
        if a1 < a2 {
            std::cmp::Ordering::Greater
        } else if a1 > a2 {
            std::cmp::Ordering::Less
        } else {
            std::cmp::Ordering::Equal
        }
    });
}

fn vaporize(mut map: HashMap<Pt2, Position>, from: &Pt2, n: usize) -> Pt2 {
    let mut targets = VecDeque::new();
    let mut vaporized = Vec::new();
    for _ in 0..n {
        if targets.len() == 0 {
            let mut next = detected_asteroids(&map, from);
            sort_angular(&mut next, from);
            assert_ne!(next.len(), 0);
            for pt in next {
                targets.push_back(*pt);
            }
        }
        let next = targets.pop_front().unwrap();
        vaporized.push(next);
        map.remove(&next).unwrap();
    }
    *vaporized.last().unwrap()
}

fn main() {
    let path = env::args().nth(1).unwrap();
    let text = fs::read_to_string(&path).unwrap();
    let map = read_map(&text);

    let pos = *best_position(&map);
    println!("{:?}", detected_asteroids(&map, &pos).len());

    let vaporized200th = vaporize(map, &pos, 200);
    println!("{:?}", vaporized200th.x * 100 + vaporized200th.y);
}
