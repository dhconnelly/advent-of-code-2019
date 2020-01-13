use std::collections::HashMap;
use std::collections::VecDeque;
use std::env;
use std::fs;

#[derive(Clone, Copy, Eq, PartialEq, Hash, Debug)]
enum Tile {
    Entrance,
    Wall,
    Passage,
    Key(u8),
    Door(u8),
}

fn read_map(s: &str) -> HashMap<geom::Point2, Tile> {
    let mut m = HashMap::new();
    let mut p = geom::ZERO2;
    for line in s.trim().lines() {
        for ch in line.bytes() {
            let t = match ch {
                b'@' => Tile::Entrance,
                b'#' => Tile::Wall,
                b'.' => Tile::Passage,
                b'a'..=b'z' => Tile::Key(ch),
                b'A'..=b'Z' => Tile::Door(ch),
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

#[derive(Copy, Clone)]
struct BfsNode {
    n: Node,
    dist: i32,
}

impl BfsNode {
    fn new(p: geom::Point2, t: Tile, dist: i32) -> BfsNode {
        BfsNode {
            n: Node { p, t },
            dist,
        }
    }
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

fn bfs(from: &geom::Point2, tiles: &HashMap<geom::Point2, Tile>) -> HashMap<Node, i32> {
    let mut es = HashMap::new();
    let mut q = VecDeque::<BfsNode>::new();
    let from_node = BfsNode::new(*from, tiles[from], 0);
    q.push_back(from_node.clone());
    let mut v = HashMap::<geom::Point2, bool>::new();
    v.insert(*from, true);
    while q.len() > 0 {
        let head = q.pop_front().unwrap();
        for nbr in neighbors(&head.n.p, &tiles, &v) {
            let nbr_t = &tiles[&nbr];
            v.insert(nbr, true);
            match nbr_t {
                Tile::Key(_) | Tile::Door(_) => {
                    es.insert(Node { p: nbr, t: *nbr_t }, head.dist + 1);
                }
                _ => (),
            }
            match nbr_t {
                Tile::Key(_) | Tile::Passage | Tile::Entrance => {
                    q.push_back(BfsNode::new(nbr, *nbr_t, head.dist + 1))
                }
                _ => (),
            }
        }
    }
    es
}

fn reachable_graph(tiles: &HashMap<geom::Point2, Tile>) -> HashMap<Node, HashMap<Node, i32>> {
    tiles
        .iter()
        .filter(|(_, t)| match t {
            Tile::Key(_) | Tile::Door(_) | Tile::Entrance => true,
            _ => false,
        })
        .map(|(p, t)| (Node { p: *p, t: *t }, bfs(p, tiles)))
        .collect()
}

fn find_entrance(g: &HashMap<Node, HashMap<Node, i32>>) -> Node {
    *g.keys().filter(|n| n.t == Tile::Entrance).next().unwrap()
}

#[derive(Hash, PartialEq, Eq, Clone)]
struct State {
    from: Node,
    have: i32,
}

impl State {
    fn new(from: Node, have: i32) -> State {
        State { from, have }
    }
}

fn is_key(n: &Node) -> bool {
    if let Tile::Key(_) = n.t {
        return true;
    }
    false
}

fn is_door(n: &Node) -> bool {
    if let Tile::Door(_) = n.t {
        return true;
    }
    false
}

fn value(n: &Node) -> u8 {
    match n.t {
        Tile::Door(x) => x,
        Tile::Key(x) => x,
        _ => panic!("no value"),
    }
}

fn door_for_key(key: &Node, g: &HashMap<Node, HashMap<Node, i32>>) -> Option<Node> {
    let val = value(key);
    g.keys()
        .copied()
        .filter(is_door)
        .filter(|nd| value(nd) == val - 32)
        .next()
}

fn remove(h: &mut HashMap<Node, HashMap<Node, i32>>, nd: &Node) {
    for (_, es) in h {
        es.remove(nd);
    }
}

fn connect_via(h: &mut HashMap<Node, HashMap<Node, i32>>, nd: &Node) {
    let mut new_edges = HashMap::<Node, HashMap<Node, i32>>::new();
    for (to1, d1) in &h[&nd] {
        for (to2, d2) in &h[&nd] {
            if to1 == to2 {
                continue;
            }
            if h[to2].contains_key(to1) {
                continue;
            }
            new_edges.entry(*to1).or_default().insert(*to2, d1 + d2);
        }
    }
    for (k1, v1) in new_edges.iter() {
        for (k2, v2) in v1.iter() {
            h.entry(*k1).or_default().insert(*k2, *v2);
        }
    }
}

fn take_key(
    g: &HashMap<Node, HashMap<Node, i32>>,
    key: &Node,
) -> HashMap<Node, HashMap<Node, i32>> {
    let mut h = g.clone();
    if let Some(door) = door_for_key(key, g) {
        connect_via(&mut h, &door);
        remove(&mut h, &door);
        h.remove(&door);
    }
    remove(&mut h, key);
    h
}

fn key_neighbors<'a>(
    from: &Node,
    g: &'a HashMap<Node, HashMap<Node, i32>>,
) -> Vec<(&'a Node, &'a i32)> {
    if let Some(edges) = g.get(from) {
        return edges.iter().filter(|(n, _)| is_key(*n)).collect();
    }
    vec![]
}

fn have_with(have: i32, item: usize) -> i32 {
    have | (1 << item)
}

fn which(key: &Node) -> usize {
    (value(key) - b'a') as usize
}

fn print_graph(g: &HashMap<Node, HashMap<Node, i32>>) {
    for (from, tos) in g {
        for (to, dist) in tos {
            println!("{:?} -> {:?} ({})", from.t, to.t, dist);
        }
    }
}

fn shortest_path_from(
    from: &Node,
    have: i32,
    remaining: usize,
    g: &HashMap<Node, HashMap<Node, i32>>,
    memo: &mut HashMap<State, i32>,
) -> Option<i32> {
    if remaining == 0 {
        return Some(0);
    }
    let state = State::new(*from, have);
    if memo.contains_key(&state) {
        return Some(memo[&state]);
    }
    //println!("{:?}, {:b}, {:?}", from, have, remaining);
    //print_graph(g);
    let mut min_dist = None;
    for (key, key_dist) in key_neighbors(from, g) {
        let g_next = take_key(g, &key);
        let have_next = have_with(have, which(key));
        let dist = match shortest_path_from(&key, have_next, remaining - 1, &g_next, memo) {
            None => continue,
            Some(x) => x + key_dist,
        };
        if min_dist == None || dist < min_dist.unwrap() {
            min_dist = Some(dist);
            memo.insert(state.clone(), dist);
        }
    }
    min_dist
}

fn shortest_path(g: &HashMap<Node, HashMap<Node, i32>>) -> i32 {
    let from = find_entrance(&g);
    let g = g.clone();
    let mut memo = HashMap::new();
    let remaining = g.keys().copied().filter(is_key).count();
    shortest_path_from(&from, 0, remaining, &g, &mut memo).unwrap()
}

fn main() {
    let path = env::args().nth(1).expect("missing input path");
    let s = fs::read_to_string(path).expect("can't read input");
    let m = read_map(&s);
    let g = reachable_graph(&m);
    println!("{}", shortest_path(&g));
}
