use geom::*;
use std::collections::HashMap;
use std::collections::HashSet;
use std::collections::VecDeque;
use std::env;
use std::error;
use std::fs;

#[derive(PartialEq)]
enum Tile {
    Passage,
    Wall,
    Empty,
    Portal(char),
}

fn read_chars(s: &str) -> HashMap<Point2, char> {
    s.lines()
        .enumerate()
        .flat_map(|(row, line)| {
            line.chars()
                .enumerate()
                .map(move |(col, ch)| (Point2::new(col as i32, row as i32), ch))
        })
        .collect()
}

type Tiles = HashMap<Point2, Tile>;

fn read_tile((pt, ch): (Point2, char)) -> Result<(Point2, Tile), String> {
    let tile = match ch {
        ' ' => Tile::Empty,
        '#' => Tile::Wall,
        '.' => Tile::Passage,
        x @ 'A'..='Z' => Tile::Portal(x),
        x => return Err(format!("bad tile at {:?}: {}", pt, x)),
    };
    Ok((pt, tile))
}

fn read_tiles(s: &str) -> Result<Tiles, String> {
    read_chars(s).into_iter().map(read_tile).collect()
}

struct Bounds<'a> {
    outer_lo: &'a Point2,
    outer_hi: &'a Point2,
    inner_lo: &'a Point2,
    inner_hi: &'a Point2,
}

fn outer_bounds(tiles: &Tiles) -> (&Point2, &Point2) {
    let mut min = tiles.keys().next().unwrap();
    let mut max = tiles.keys().next().unwrap();
    for (pt, _) in tiles.iter().filter(|(_, tile)| *tile == &Tile::Wall) {
        if pt.x < min.x || pt.y < min.y {
            min = pt;
        }
        if pt.x > max.x || pt.y > max.y {
            max = pt;
        }
    }
    (min, max)
}

fn in_bounds(pt: &Point2, lo: &Point2, hi: &Point2) -> bool {
    pt.x > lo.x && pt.y > lo.y && pt.x < hi.x && pt.y < hi.y
}

fn inner_bounds<'a>(
    tiles: &'a Tiles,
    outer_lo: &Point2,
    outer_hi: &Point2,
) -> (&'a Point2, &'a Point2) {
    let mut min = tiles.keys().next().unwrap();
    let mut max = tiles.keys().next().unwrap();
    for (pt, _) in tiles.iter().filter(|(pt, tile)| {
        in_bounds(pt, outer_lo, outer_hi) && *tile == &Tile::Empty
    }) {
        if pt.x < min.x || pt.y < min.y {
            min = pt;
        }
        if pt.x > max.x || pt.y > max.y {
            max = pt;
        }
    }
    (min, max)
}

fn maze_bounds(tiles: &Tiles) -> Bounds {
    let (outer_lo, outer_hi) = outer_bounds(tiles);
    let (inner_lo, inner_hi) = inner_bounds(tiles, outer_lo, outer_hi);
    Bounds {
        outer_lo,
        outer_hi,
        inner_lo,
        inner_hi,
    }
}

fn manhattan_nbrs<'a>(pt: &'a Point2, tiles: &'a Tiles) -> HashSet<&'a Point2> {
    pt.manhattan_neighbors()
        .iter()
        .filter_map(|nbr| match tiles.get_key_value(nbr) {
            Some((nbr, Tile::Passage)) => Some(nbr),
            Some(_) => None,
            None => None,
        })
        .collect()
}

fn all_manhattan_nbrs(tiles: &Tiles) -> HashMap<&Point2, HashSet<&Point2>> {
    tiles
        .iter()
        .filter(|(_, tile)| *tile == &Tile::Passage)
        .map(|(pt, _)| (pt, manhattan_nbrs(pt, tiles)))
        .collect()
}

type Portal = (char, char);

fn portal_from(p: &Point2, dir: Direction, tiles: &Tiles) -> Option<Portal> {
    let r = p.go(dir);
    let s = r.go(dir);
    match (tiles.get(&r).unwrap(), tiles.get(&s).unwrap()) {
        (Tile::Portal(c1), Tile::Portal(c2)) => Some((*c1, *c2)),
        _ => None,
    }
}

fn portals_along<T: IntoIterator<Item = Point2>>(
    pts: T,
    dir: Direction,
    invert: bool,
    tiles: &Tiles,
) -> impl Iterator<Item = (Portal, &Point2)> {
    pts.into_iter().filter_map(move |p| {
        let p = tiles.get_key_value(&p).unwrap().0;
        let (c1, c2) = portal_from(p, dir, tiles)?;
        let (c1, c2) = if invert { (c2, c1) } else { (c1, c2) };
        Some(((c1, c2), p))
    })
}

fn all_portals<'a>(
    tiles: &'a Tiles,
    Bounds {
        outer_lo,
        outer_hi,
        inner_lo,
        inner_hi,
    }: &Bounds,
) -> HashMap<Portal, Vec<&'a Point2>> {
    let pt = Point2::new;
    let outer_top = (outer_lo.x..=outer_hi.x).map(|x| pt(x, outer_lo.y));
    let outer_bottom = (outer_lo.x..=outer_hi.x).map(|x| pt(x, outer_hi.y));
    let outer_left = (outer_lo.y..=outer_hi.y).map(|y| pt(outer_lo.x, y));
    let outer_right = (outer_lo.y..=outer_hi.y).map(|y| pt(outer_hi.x, y));
    let inner_top = (inner_lo.x..=inner_hi.x).map(|x| pt(x, inner_lo.y - 1));
    let inner_bottom = (inner_lo.x..=inner_hi.x).map(|x| pt(x, inner_hi.y + 1));
    let inner_left = (inner_lo.y..=inner_hi.y).map(|y| pt(inner_lo.x - 1, y));
    let inner_right = (inner_lo.y..=inner_hi.y).map(|y| pt(inner_hi.x + 1, y));

    let mut portals = HashMap::new();
    for (portal, p) in portals_along(outer_top, Direction::Down, true, tiles)
        .chain(portals_along(outer_bottom, Direction::Up, false, tiles))
        .chain(portals_along(outer_left, Direction::Left, true, tiles))
        .chain(portals_along(outer_right, Direction::Right, false, tiles))
        .chain(portals_along(inner_top, Direction::Up, false, tiles))
        .chain(portals_along(inner_bottom, Direction::Down, true, tiles))
        .chain(portals_along(inner_left, Direction::Right, false, tiles))
        .chain(portals_along(inner_right, Direction::Left, true, tiles))
    {
        portals.entry(portal).or_insert(vec![]).push(p);
    }
    portals
}

fn link_portals<'a>(
    nbrs: &mut HashMap<&'a Point2, HashSet<&'a Point2>>,
    portals: &HashMap<Portal, Vec<&'a Point2>>,
) {
    for pq in portals.values() {
        for p in pq.iter() {
            for q in pq.iter() {
                if p != q {
                    nbrs.entry(q).or_default().insert(p);
                    nbrs.entry(p).or_default().insert(q);
                }
            }
        }
    }
}

type Graph<'a> = HashMap<&'a Point2, HashSet<&'a Point2>>;

struct Maze<'a> {
    begin: &'a Point2,
    end: &'a Point2,
    nbrs: Graph<'a>,
    bounds: Bounds<'a>,
}

fn read_maze(tiles: &Tiles) -> Maze {
    let mut nbrs = all_manhattan_nbrs(tiles);
    let bounds = maze_bounds(tiles);
    let portals = all_portals(tiles, &bounds);
    link_portals(&mut nbrs, &portals);
    let begin = portals[&('A', 'A')][0];
    let end = portals[&('Z', 'Z')][0];
    Maze {
        begin,
        end,
        nbrs,
        bounds,
    }
}

#[derive(Debug, Clone, Copy, Eq, PartialEq, Hash)]
struct RecPt2<'a> {
    p: &'a Point2,
    lvl: usize,
}

impl<'a> RecPt2<'a> {
    fn new(p: &'a Point2, lvl: usize) -> RecPt2<'a> {
        RecPt2 { p, lvl }
    }
}

fn dist<'a>(
    maze: &'a Maze,
    nbrs: impl Fn(&RecPt2, &'a Maze) -> Vec<RecPt2<'a>>,
    eq: impl Fn(&RecPt2, &RecPt2) -> bool,
) -> Option<usize> {
    let src = RecPt2::new(maze.begin, 0);
    let dst = RecPt2::new(maze.end, 0);
    let mut q = VecDeque::new();
    let mut v = HashSet::new();
    q.push_back((src, 0));
    v.insert(src);
    while let Some((p, d)) = q.pop_front() {
        let ns = nbrs(&p, &maze);
        for nbr in ns {
            if eq(&nbr, &dst) {
                return Some(d + 1);
            }
            if v.contains(&nbr) {
                continue;
            }
            v.insert(nbr);
            q.push_back((nbr, d + 1));
        }
    }
    None
}

fn flat_eq(p: &RecPt2, q: &RecPt2) -> bool {
    p.p == q.p
}

fn flat_nbrs<'a>(p: &RecPt2, maze: &'a Maze) -> Vec<RecPt2<'a>> {
    maze.nbrs
        .get(p.p)
        .unwrap()
        .iter()
        .map(|p| RecPt2 { p, lvl: 0 })
        .collect()
}

fn flat_dist(maze: &Maze) -> Option<usize> {
    dist(&maze, flat_nbrs, flat_eq)
}

fn rec_eq(p: &RecPt2, q: &RecPt2) -> bool {
    p == q
}

fn is_outer_edge(p: &RecPt2, maze: &Maze) -> bool {
    p.p.x == maze.bounds.outer_lo.x
        || p.p.x == maze.bounds.outer_hi.x
        || p.p.y == maze.bounds.outer_lo.y
        || p.p.y == maze.bounds.outer_hi.y
}

fn is_inner_edge(p: &RecPt2, maze: &Maze) -> bool {
    p.p.x + 1 == maze.bounds.inner_lo.x
        || p.p.x - 1 == maze.bounds.inner_hi.x
        || p.p.y + 1 == maze.bounds.inner_lo.y
        || p.p.y - 1 == maze.bounds.inner_hi.y
}

fn rec_nbrs<'a>(p: &RecPt2, maze: &'a Maze) -> Vec<RecPt2<'a>> {
    let mut nbrs = Vec::new();
    for nbr in maze.nbrs.get(p.p).unwrap() {
        if nbr.manhattan_dist(p.p) == 1 {
            // non-portal neighbors are always accessible.
            nbrs.push(RecPt2::new(nbr, p.lvl));
        } else if is_inner_edge(p, maze) {
            // portal neighbors from the inner edge are always accessible,
            // but at a greater recursion depth.
            nbrs.push(RecPt2::new(nbr, p.lvl + 1));
        } else if is_outer_edge(p, maze) && p.lvl > 0 {
            // portal neighbors from the outer edge are only accessible from
            // non-zero recursion depth. traversing them goes up a level.
            nbrs.push(RecPt2::new(nbr, p.lvl - 1));
        } else {
            // portal neighbors from the outer edge at level zero are closed.
            assert!(is_outer_edge(p, maze) && p.lvl == 0);
        }
    }
    nbrs
}

fn rec_dist(maze: &Maze) -> Option<usize> {
    dist(&maze, rec_nbrs, rec_eq)
}

fn main() -> Result<(), Box<dyn error::Error>> {
    let path = env::args().nth(1).ok_or("Usage: day20 <filename>")?;
    let text = fs::read_to_string(&path)?;
    let tiles = read_tiles(&text)?;
    let maze = read_maze(&tiles);
    println!("{}", flat_dist(&maze).unwrap());
    println!("{}", rec_dist(&maze).unwrap());
    Ok(())
}
