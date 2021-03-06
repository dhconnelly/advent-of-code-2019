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
    Key(u8),
    Door(u8),
}

#[derive(Clone, Debug, PartialEq, Eq, Hash)]
struct Node {
    p: geom::Point2,
    t: Tile,
}

impl Node {
    fn key_value(&self) -> u8 {
        if let Tile::Key(key) = self.t {
            return key;
        }
        panic!("not a key");
    }
}

#[derive(Clone)]
struct BfsNode<'a> {
    node: &'a Node,
    dist: i32,
}

impl BfsNode<'_> {
    fn new(node: &Node, dist: i32) -> BfsNode {
        BfsNode { node, dist }
    }
}

#[derive(Clone, Debug, PartialEq, Eq, Hash)]
struct KeySet {
    keys: usize,
}

impl KeySet {
    fn new() -> KeySet {
        KeySet { keys: 0 }
    }

    fn with(&self, key: u8) -> KeySet {
        let mut ks = self.clone();
        ks.add(key);
        ks
    }

    fn add(&mut self, key: u8) {
        self.keys |= 1 << (key - b'a') as usize;
    }

    fn unlocks_door(&self, door: u8) -> bool {
        self.has(door + 32u8)
    }

    fn has(&self, key: u8) -> bool {
        (self.keys & (1 << (key - b'a') as usize)) > 0
    }
}

impl fmt::Display for KeySet {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{:b}", self.keys)
    }
}

#[derive(Debug)]
struct Map {
    tiles: HashMap<geom::Point2, Node>,
}

impl Map {
    fn new(s: &str) -> Map {
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
                m.insert(p, Node { p, t });
                p.x += 1;
            }
            p.x = 0;
            p.y += 1;
        }
        Map { tiles: m }
    }

    fn entrance(&self) -> Option<&Node> {
        self.tiles.values().filter(|n| n.t == Tile::Entrance).next()
    }

    fn passable(&self, node: &Node, keys: &KeySet) -> bool {
        match node.t {
            Tile::Entrance | Tile::Passage | Tile::Key(_) => true,
            Tile::Door(door) => keys.unlocks_door(door),
            Tile::Wall => false,
        }
    }

    fn keys(&self) -> impl Iterator<Item = &Node> {
        self.tiles.values().filter(|n| {
            if let Tile::Key(_) = n.t {
                return true;
            }
            false
        })
    }
}

struct Explorer<'a> {
    map: &'a Map,
    q: VecDeque<BfsNode<'a>>,
    v: HashMap<&'a geom::Point2, bool>,
}

impl<'a> Explorer<'a> {
    fn new(map: &Map) -> Explorer {
        Explorer {
            map: map,
            q: VecDeque::new(),
            v: HashMap::new(),
        }
    }

    fn reachable_keys(&mut self, from: &'a Node, with_keys: &KeySet) -> Vec<BfsNode<'a>> {
        let mut keys = Vec::new();
        let q = &mut self.q;
        q.clear();
        q.push_back(BfsNode::new(from, 0));
        let v = &mut self.v;
        v.clear();
        v.insert(&from.p, true);
        while !q.is_empty() {
            let front = q.pop_front().unwrap();
            if let Tile::Key(key) = front.node.t {
                if !with_keys.has(key) {
                    keys.push(front.clone());
                }
            }
            for nbr in &front.node.p.manhattan_neighbors() {
                if let Some(nbr) = self.map.tiles.get(nbr) {
                    if self.map.passable(nbr, with_keys) && !v.contains_key(&nbr.p) {
                        q.push_back(BfsNode::new(nbr, front.dist + 1));
                        v.insert(&front.node.p, true);
                    }
                }
            }
        }
        keys
    }

    fn shortest_path_with(
        &mut self,
        from: &'a Node,
        keys: &KeySet,
        remaining: i32,
        memo: &mut HashMap<MemoKey, Option<i32>>,
    ) -> Option<i32> {
        let mk = MemoKey::new(from, keys);
        if memo.contains_key(&mk) {
            return memo[&mk];
        }
        if remaining == 0 {
            return Some(0);
        }
        let mut min_dist = None;
        let reachable = self.reachable_keys(from, keys);
        for key in &reachable {
            let node = key.node;
            match self.shortest_path_with(node, &keys.with(node.key_value()), remaining - 1, memo) {
                None => continue,
                Some(dist) => {
                    let dist = dist + key.dist;
                    let min = *min_dist.get_or_insert(dist);
                    min_dist.replace(min.min(dist));
                }
            }
        }
        memo.insert(mk, min_dist);
        min_dist
    }
}

#[derive(Hash, PartialEq, Eq, Debug)]
struct MemoKey {
    from: Node,
    keys: KeySet,
}

impl MemoKey {
    fn new(from: &Node, keys: &KeySet) -> MemoKey {
        MemoKey {
            from: from.clone(),
            keys: keys.clone(),
        }
    }
}

fn shortest_path(map: &Map) -> i32 {
    let mut explorer = Explorer::new(map);
    explorer
        .shortest_path_with(
            map.entrance().unwrap(),
            &KeySet::new(),
            map.keys().count() as i32,
            &mut HashMap::new(),
        )
        .unwrap()
}

fn main() {
    let path = env::args().nth(1).expect("missing input path");
    let s = fs::read_to_string(path).expect("can't read input");
    let m = Map::new(&s);
    println!("{}", shortest_path(&m));
}
