use std::env;
use std::error::Error;
use std::fs;

fn valid1(mut x: i32) -> bool {
    let mut prev = x % 10;
    let mut chain = 1;
    let mut two_chain = false;
    x /= 10;
    while x > 0 {
        let y = x % 10;
        if y > prev {
            return false;
        }
        if y == prev {
            chain += 1;
        } else {
            chain = 1;
        }
        if chain == 2 {
            two_chain = true;
        }
        prev = y;
        x /= 10;
    }
    two_chain
}

fn valid2(mut x: i32) -> bool {
    let mut prev = x % 10;
    let mut chain = 1;
    let mut two_chain = false;
    x /= 10;
    while x > 0 {
        let y = x % 10;
        if y > prev {
            return false;
        }
        if y == prev {
            chain += 1;
        } else {
            if chain == 2 {
                two_chain = true;
            }
            chain = 1;
        }
        prev = y;
        x /= 10;
    }
    two_chain || chain == 2
}

fn count_valid<T: Fn(i32) -> bool>(from: i32, to: i32, valid: T) -> usize {
    (from..to + 1).filter(|x| valid(*x)).count()
}

fn main() -> Result<(), Box<dyn Error>> {
    let path = env::args().nth(1).ok_or("missing input path")?;
    let toks = fs::read_to_string(path)?
        .trim()
        .split('-')
        .map(|x| x.parse::<i32>())
        .collect::<Result<Vec<_>, _>>()?;
    let (from, to) = (toks[0], toks[1]);
    println!("{}", count_valid(from, to, valid1));
    println!("{}", count_valid(from, to, valid2));
    Ok(())
}
