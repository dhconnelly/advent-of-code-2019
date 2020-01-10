use std::collections::HashMap;
use std::collections::HashSet;
use std::error::Error;
use std::iter;

fn read_orbits(input: &str) -> HashMap<&str, &str> {
    let mut orbits = HashMap::new();
    for line in input.trim().lines() {
        let mut toks = line.splitn(2, ')');
        let (left, right) = (toks.next().unwrap(), toks.next().unwrap());
        orbits.insert(right, left);
    }
    orbits
}

fn count_orbits<'a>(
    from: &'a str,
    orbits: &HashMap<&'a str, &'a str>,
    counts: &mut HashMap<&'a str, i32>,
) -> i32 {
    if let Some(count) = counts.get(from) {
        return *count;
    }
    let count = match orbits.get(from) {
        None => 0,
        Some(to) => count_orbits(to, orbits, counts) + 1,
    };
    counts.insert(from, count);
    count
}

fn count_all_orbits(orbits: &HashMap<&str, &str>) -> i32 {
    let mut counts = HashMap::new();
    orbits
        .keys()
        .map(|from| count_orbits(from, orbits, &mut counts))
        .sum()
}

fn orbit_path<'a>(from: &'a str, orbits: &HashMap<&str, &'a str>) -> Vec<&'a str> {
    match orbits.get(from) {
        Some(to) => iter::once(from).chain(orbit_path(to, orbits)).collect(),
        None => vec![from],
    }
}

fn dist(from: &str, to: &str, orbits: &HashMap<&str, &str>) -> usize {
    let orbit = orbit_path(from, orbits);
    let set: HashSet<&str> = orbit_path(from, orbits).iter().copied().collect();
    let (i, parent) = orbit_path(to, orbits)
        .iter()
        .copied()
        .enumerate()
        .find(|(_, node)| set.contains(node))
        .unwrap();
    let j = orbit.iter().position(|node| *node == parent).unwrap();
    i + j - 2
}

fn main() -> Result<(), Box<dyn Error>> {
    let path = std::env::args().nth(1).ok_or("missing path")?;
    let input = std::fs::read_to_string(path)?;
    let orbits = read_orbits(&input);
    println!("{}", count_all_orbits(&orbits));
    println!("{}", dist("YOU", "SAN", &orbits));
    Ok(())
}
