use std::env;
use std::error;
use std::fs;

fn run(program: &intcode::Program, noun: i64, verb: i64) -> Result<i64, String> {
    let mut machine = intcode::Machine::new(program);
    machine.set(1, noun);
    machine.set(2, verb);
    assert_eq!(machine.run()?, intcode::State::Halted);
    Ok(machine.get(0))
}

fn find_noun_verb(program: &intcode::Program, target: i64) -> Result<(i64, i64), String> {
    for noun in 0..100 {
        for verb in 0..100 {
            if run(&program, noun, verb)? == target {
                return Ok((noun, verb));
            }
        }
    }
    Err(format!("no noun and verb matching {}", target))
}

fn main() -> Result<(), Box<dyn error::Error>> {
    let path = env::args().nth(1).ok_or("missing input path")?;
    let text = fs::read_to_string(&path)?;
    let program = intcode::Program::new(&text)?;
    println!("{}", run(&program, 12, 2)?);
    let (noun, verb) = find_noun_verb(&program, 19690720)?;
    println!("{}", 100 * noun + verb);
    Ok(())
}
