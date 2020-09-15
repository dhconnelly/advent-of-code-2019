use std::collections::HashMap;
use std::env;
use std::fs;

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
struct Pt2 {
    x: usize,
    y: usize,
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
            map.insert(Pt2 { x, y }, Position::of(ch));
        }
    }
    map
}

fn asteroids(map: &HashMap<Pt2, Position>) -> Vec<&Pt2> {
    map.iter().filter(|(_, pos)| **pos == Position::Asteroid).map(|(pt, _)| pt).collect()
}

const EPS: f64 = 0.00001;

fn angle(from: &Pt2, to: &Pt2) -> f64 {
    (to.y as f64 - from.y as f64).atan2(to.x as f64 - from.x as f64)
}

fn detected_asteroids<'a>(map: &'a HashMap<Pt2, Position>, from: &Pt2) -> Vec<&'a Pt2> {
    let mut detected: Vec<(&'a Pt2, f64)> = Vec::new();
    for pt in asteroids(map) {
        if pt == from {
            continue;
        }
        let a = angle(from, pt);
        if detected.iter().filter(|(_, a2)| (a - *a2).abs() < EPS).count() == 0 {
            detected.push((pt, a));
        }
    }
    detected.into_iter().map(|(pt, _)| pt).collect()
}

fn best_position(map: &HashMap<Pt2, Position>) -> &Pt2 {
    asteroids(map).into_iter()
        .map(|pt| (pt, detected_asteroids(map, pt).len()))
        .max_by(|(_, n1), (_, n2)| n1.cmp(n2))
        .unwrap()
        .0
}

fn main() {
    let path = env::args().nth(1).unwrap();
    let text = fs::read_to_string(&path).unwrap();
    let map = read_map(&text);
    println!("{:?}", detected_asteroids(&map, best_position(&map)).len());
}
