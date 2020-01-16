use std::env;
use std::error;
use std::fs;

fn run(prog: &intcode::Program, input: i64) -> Result<i64, String> {
    let mut machine = intcode::Machine::new(prog);
    let mut output = Vec::new();
    loop {
        let state = machine.run()?;
        match state {
            intcode::State::Running => (),
            intcode::State::Reading => machine.write(input),
            intcode::State::Writing => output.push(machine.read()),
            intcode::State::Halted => break,
        }
    }
    if output.len() != 1 {
        Err(format!("bad output: {:?}", output))
    } else {
        Ok(output[0])
    }
}

fn main() -> Result<(), Box<dyn error::Error>> {
    let path = env::args().nth(1).ok_or("missing input path")?;
    let text = fs::read_to_string(&path)?;
    let prog = intcode::Program::new(&text)?;
    println!("{}", run(&prog, 1)?);
    println!("{}", run(&prog, 2)?);
    Ok(())
}
