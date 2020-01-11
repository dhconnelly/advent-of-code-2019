use std::collections::HashSet;
use std::env;
use std::error::Error;
use std::fs;

#[derive(Copy, Clone, Debug)]
struct Step {
    dir: geom::Direction,
    dist: i32,
}

fn read_path(line: &str) -> Result<Vec<Step>, String> {
    line.split(',')
        .map(|tok| {
            let dir = match &tok[..1] {
                "R" => Ok(geom::Direction::Right),
                "L" => Ok(geom::Direction::Left),
                "U" => Ok(geom::Direction::Up),
                "D" => Ok(geom::Direction::Down),
                _ => Err(format!("bad step {} in line {}", tok, line)),
            }?;
            let dist = tok[1..]
                .parse::<i32>()
                .map_err(|_| format!("bad step {} in line {}", tok, line))?;
            Ok(Step { dir, dist })
        })
        .collect()
}

fn read_paths(input: &str) -> Result<Vec<Vec<Step>>, String> {
    input.trim().lines().take(2).map(|l| read_path(l)).collect()
}

fn read_wire(path: &Vec<Step>) -> Vec<geom::Point2> {
    let mut p = geom::ZERO2;
    let mut w = Vec::new();
    for step in path {
        for i in 0..step.dist {
            p = p.go(step.dir);
            w.push(p);
        }
    }
    w
}

fn read_wires(paths: &Vec<Vec<Step>>) -> Vec<Vec<geom::Point2>> {
    paths.iter().map(read_wire).collect()
}

fn closest_intersect(p1: &Vec<geom::Point2>, p2: &Vec<geom::Point2>) -> i32 {
    let set: HashSet<&geom::Point2> = p1.iter().collect();
    let mut min_dist = std::i32::MAX;
    for p in p2 {
        let dist = p.manhattan_norm();
        if set.contains(p) {
            if dist < min_dist {
                min_dist = dist;
            }
        }
    }
    min_dist
}

fn fastest_intersect(p1: &Vec<geom::Point2>, p2: &Vec<geom::Point2>) -> i32 {
    let set: HashSet<&geom::Point2> = p1.iter().collect();
    let mut min_dist = std::i32::MAX;
    for (i, p) in p2.iter().enumerate() {
        if set.contains(p) {
            let j = p1.iter().position(|q| q == p).unwrap();
            let dist = (i + j) as i32 + 2;
            if dist < min_dist {
                min_dist = dist;
            }
        }
    }
    min_dist
}

fn main() -> Result<(), Box<dyn Error>> {
    let path = env::args().nth(1).ok_or("missing input path")?;
    let input = fs::read_to_string(path)?;
    let paths = read_paths(&input)?;
    let wires = read_wires(&paths);
    println!("{}", closest_intersect(&wires[0], &wires[1]));
    println!("{}", fastest_intersect(&wires[0], &wires[1]));
    Ok(())
}
