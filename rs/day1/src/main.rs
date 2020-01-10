use std::env;
use std::error::Error;
use std::fs;

fn fuel(mass: i32) -> i32 {
    (mass / 3) - 2
}

fn rec_fuel(mass: i32) -> i32 {
    let f = fuel(mass);
    if f <= 0 {
        0
    } else {
        f + rec_fuel(f)
    }
}

fn main() -> Result<(), Box<dyn Error>> {
    let path = env::args().nth(1).ok_or("missing input path")?;
    let input: Vec<i32> = fs::read_to_string(path)?
        .lines()
        .map(|l| l.parse::<i32>().unwrap())
        .collect();
    let part1: i32 = input.iter().copied().map(fuel).sum();
    println!("{}", part1);
    let part2: i32 = input.iter().copied().map(rec_fuel).sum();
    println!("{}", part2);
    Ok(())
}
