use geom::*;
use std::collections::HashMap;
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

#[derive(Debug)]
struct Maze {
    nbrs: HashMap<Point2, Vec<Point2>>,
    begin: Point2,
    end: Point2,
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

fn manhattan_nbrs<'a>(pt: &'a Point2, tiles: &'a Tiles) -> Vec<&'a Point2> {
    pt.manhattan_neighbors()
        .iter()
        .filter_map(|nbr| match tiles.get_key_value(nbr) {
            Some((nbr, Tile::Passage)) => Some(nbr),
            Some(_) => None,
            None => None,
        })
        .collect()
}

fn all_manhattan_nbrs(tiles: &Tiles) -> HashMap<&Point2, Vec<&Point2>> {
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
    }: Bounds<'a>,
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

fn read_maze(tiles: &Tiles) -> Result<Maze, String> {
    let nbrs = all_manhattan_nbrs(tiles);
    let bounds = maze_bounds(tiles);
    let portals = all_portals(tiles, bounds);
    for (portal, nbrs) in portals.iter() {
        println!("{:?} {:?}", portal, nbrs);
    }
    // link across portals
    // find begin and end
    Err("not implemented".to_string())
}

fn main() -> Result<(), Box<dyn error::Error>> {
    let path = env::args().nth(1).ok_or("Usage: day20 <filename>")?;
    let text = fs::read_to_string(&path)?;
    let tiles = read_tiles(&text)?;
    let maze = read_maze(&tiles)?;
    println!("{:?}", maze);
    Ok(())
}
