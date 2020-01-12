use std::collections::HashMap;
use std::collections::VecDeque;
use std::env;
use std::fmt;
use std::fs;

#[derive(Clone, Copy, Eq, PartialEq, Hash, Debug)]
enum Tile {
    Entrance,
    Wall,
    Passage,
    Key(char),
    Door(char),
}

fn read_map(s: &str) -> HashMap<geom::Point2, Tile> {
    let mut m = HashMap::new();
    let mut p = geom::ZERO2;
    for line in s.trim().lines() {
        for ch in line.chars() {
            let t = match ch {
                '@' => Tile::Entrance,
                '#' => Tile::Wall,
                '.' => Tile::Passage,
                'a'..='z' => Tile::Key(ch),
                'A'..='Z' => Tile::Door(ch),
                _ => panic!(format!("bad tile: {}", ch)),
            };
            m.insert(p, t);
            p.x += 1;
        }
        p.x = 0;
        p.y += 1;
    }
    m
}

#[derive(Copy, Clone, PartialEq, Eq, Hash, Debug)]
struct Node {
    p: geom::Point2,
    t: Tile,
}

#[derive(Copy, Clone, PartialEq, Eq, Debug)]
struct Edge {
    from: Node,
    to: Node,
    dist: i32,
}

struct BfsNode {
    p: geom::Point2,
    t: Tile,
    dist: i32,
}

fn neighbors(
    p: &geom::Point2,
    tiles: &HashMap<geom::Point2, Tile>,
    v: &HashMap<geom::Point2, bool>,
) -> Vec<geom::Point2> {
    p.manhattan_neighbors()
        .iter()
        .copied()
        .filter(|q| tiles.get(q).unwrap_or(&Tile::Wall) != &Tile::Wall)
        .filter(|q| !v.get(q).unwrap_or(&false))
        .collect()
}

fn bfs(from: &geom::Point2, tiles: &HashMap<geom::Point2, Tile>) -> Vec<Edge> {
    let mut es = Vec::new();
    let mut q = VecDeque::<BfsNode>::new();
    q.push_back(BfsNode {
        p: *from,
        t: tiles[from],
        dist: 0,
    });
    let mut v = HashMap::<geom::Point2, bool>::new();
    v.insert(*from, true);
    while q.len() > 0 {
        let head = q.pop_front().unwrap();
        let nbrs = head.p.manhattan_neighbors();
        for nbr in neighbors(&head.p, &tiles, &v) {
            let nbr_t = &tiles[&nbr];
            v.insert(nbr, true);
            match nbr_t {
                Tile::Key(ch) | Tile::Door(ch) => es.push(Edge {
                    from: Node {
                        p: *from,
                        t: tiles[from],
                    },
                    to: Node { p: nbr, t: *nbr_t },
                    dist: head.dist + 1,
                }),
                _ => (),
            }
            match nbr_t {
                Tile::Key(_) | Tile::Passage | Tile::Entrance => q.push_back(BfsNode {
                    p: nbr,
                    t: *nbr_t,
                    dist: head.dist + 1,
                }),
                _ => (),
            }
        }
    }
    es
}

fn reachable_graph(tiles: &HashMap<geom::Point2, Tile>) -> HashMap<Node, Vec<Edge>> {
    tiles
        .iter()
        .filter(|(p, t)| match t {
            Tile::Key(_) => true,
            Tile::Door(_) => true,
            _ => false,
        })
        .map(|(p, t)| (Node { p: *p, t: *t }, bfs(p, tiles)))
        .collect()
}

fn main() {
    let path = env::args().nth(1).expect("missing input path");
    let s = fs::read_to_string(path).expect("can't read input");
    let m = read_map(&s);
    let g = reachable_graph(&m);
    println!("{:?}", g);
}
