use std::collections::HashSet;
use std::collections::HashMap;
use std::env;
use std::error::Error;
use std::fmt;
use std::fs;

#[derive(Debug, PartialEq, Eq, Hash, Clone, Copy)]
struct Pos {
    row: usize,
    col: usize,
}

impl Pos {
    fn new(row: usize, col: usize) -> Pos {
        Pos { row, col }
    }
}

#[derive(Debug, PartialEq, Eq, Clone, Copy)]
enum State {
    Alive,
    Dead,
}

impl State {
    fn of(ch: char) -> Result<State, String> {
        match ch {
            '.' => Ok(State::Dead),
            '#' => Ok(State::Alive),
            ch => Err(format!("bad state: {}", ch)),
        }
    }
}

#[derive(Debug, Clone)]
struct Grid {
    states: HashMap<Pos, State>,
}

impl fmt::Display for Grid {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        for row in 0..5 {
            for col in 0..5 {
                if self.state(&Pos::new(row, col)) == State::Alive {
                    write!(f, "#")?
                } else {
                    write!(f, ".")?
                }
            }
            writeln!(f)?
        }
        Ok(())
    }
}

impl Grid {
    fn code(&self) -> i32 {
        self.states
            .iter()
            .map(|(Pos { row, col }, state)| match state {
                State::Alive => 1 << (row * 5 + col),
                State::Dead => 0,
            })
            .sum()
    }

    fn state(&self, pos: &Pos) -> State {
        *self.states.get(pos).unwrap_or(&State::Dead)
    }

    fn neighbors(&self, pos: &Pos) -> usize {
        let mut nbrs = 0;
        if self.state(&Pos::new(pos.row + 1, pos.col)) == State::Alive {
            nbrs = nbrs + 1;
        }
        if self.state(&Pos::new(pos.row, pos.col + 1)) == State::Alive {
            nbrs = nbrs + 1;
        }
        if pos.row > 0 && self.state(&Pos::new(pos.row - 1, pos.col)) == State::Alive {
            nbrs = nbrs + 1;
        }
        if pos.col > 0 && self.state(&Pos::new(pos.row, pos.col - 1)) == State::Alive {
            nbrs = nbrs + 1;
        }
        nbrs
    }

    fn step(&self) -> Grid {
        Grid {
            states: self
                .states
                .iter()
                .map(|(pos, state)| {
                    (
                        *pos,
                        match (state, self.neighbors(pos)) {
                            (State::Alive, 1) => State::Alive,
                            (State::Dead, 1) | (State::Dead, 2) => State::Alive,
                            _ => State::Dead,
                        },
                    )
                })
                .collect(),
        }
    }
}

fn read_grid(text: &str) -> Result<Grid, String> {
    let mut states = HashMap::new();
    for (row, line) in text.lines().enumerate() {
        for (col, ch) in line.chars().enumerate() {
            states.insert(Pos::new(row, col), State::of(ch)?);
        }
    }
    Ok(Grid { states })
}

fn repeated_code(grid: &Grid) -> i32 {
    let mut codes: HashSet<i32> = HashSet::new();
    let mut current = grid.clone();
    loop {
        let code = current.code();
        if codes.contains(&code) {
            return code;
        }
        codes.insert(code);
        current = current.step();
    }
}

fn main() -> Result<(), Box<dyn Error>> {
    let path = env::args().nth(1).ok_or("must specify input path")?;
    let text = fs::read_to_string(&path)?;
    let grid = read_grid(&text)?;
    println!("{}", repeated_code(&grid));
    Ok(())
}
