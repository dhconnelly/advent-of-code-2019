use std::env;
use std::error;
use std::fs;

fn run(prog: &intcode::Program, input: i64) -> Result<i64, String> {
    let mut machine = intcode::Machine::new(&prog);
    let mut state = machine.run()?;
    let mut output = 0;
    loop {
        if state == intcode::State::Reading {
            machine.write(input);
        } else if state == intcode::State::Writing {
            output = machine.read();
        }
        state = machine.run()?;
        if state == intcode::State::Halted {
            break;
        }
        if output != 0 {
            return Err(format!("diagnostic failed: {}", output));
        }
    }
    Ok(output)
}

fn main() -> Result<(), Box<dyn error::Error>> {
    let path = env::args().nth(1).ok_or("missing input path")?;
    let text = fs::read_to_string(&path)?;
    let prog = intcode::Program::new(&text)?;
    println!("{}", run(&prog, 1)?);
    Ok(())
}
