use std::collections::HashMap;
use std::error::Error;

#[derive(Clone, Debug)]
pub struct Program(Vec<i64>);

impl Program {
    pub fn new(text: &str) -> Result<Program, Box<dyn Error>> {
        let b: Result<Vec<_>, _> = text
            .trim()
            .split(",")
            .map(|tok| tok.parse::<i64>())
            .collect();
        Ok(Program(
            b.map_err(|e| format!("can't parse program: {}", e))?,
        ))
    }

    fn to_memory(&self) -> HashMap<i64, i64> {
        self.0
            .iter()
            .copied()
            .enumerate()
            .map(|(i, x)| (i as i64, x))
            .collect()
    }
}

#[derive(Clone)]
enum Opcode {
    Add,
    Mul,
    Halt,
}

impl Opcode {
    fn new(x: i64) -> Result<Opcode, String> {
        match x {
            1 => Ok(Opcode::Add),
            2 => Ok(Opcode::Mul),
            99 => Ok(Opcode::Halt),
            _ => Err(format!("unknown opcode: {}", x)),
        }
    }
}

#[derive(Debug)]
enum Instruction {
    Add(i64, i64, i64),
    Mul(i64, i64, i64),
    Halt,
}

#[derive(Debug, PartialEq, Copy, Clone)]
pub enum State {
    Running,
    Halted,
}

pub struct Machine {
    state: State,
    mem: HashMap<i64, i64>,
    pc: i64,
}

impl Machine {
    pub fn new(program: &Program) -> Machine {
        Machine {
            state: State::Running,
            mem: program.to_memory(),
            pc: 0,
        }
    }

    pub fn get(&self, i: i64) -> i64 {
        *self.mem.get(&i).unwrap_or(&0)
    }

    pub fn set(&mut self, i: i64, x: i64) {
        self.mem.insert(i, x);
    }

    fn get_instr(&self) -> Result<Instruction, String> {
        let arg0 = self.get(self.pc + 0);
        let arg1 = self.get(self.get(self.pc + 1));
        let arg2 = self.get(self.get(self.pc + 2));
        let arg3 = self.get(self.pc + 3);
        let instr = match Opcode::new(arg0)? {
            Opcode::Add => Instruction::Add(arg1, arg2, arg3),
            Opcode::Mul => Instruction::Mul(arg1, arg2, arg3),
            Opcode::Halt => Instruction::Halt,
        };
        Ok(instr)
    }

    fn exec(&mut self, instr: &Instruction) {
        match instr {
            Instruction::Add(x, y, z) => {
                self.set(*z, x + y);
                self.state = State::Running;
                self.pc += 4;
            }
            Instruction::Mul(x, y, z) => {
                self.set(*z, x * y);
                self.state = State::Running;
                self.pc += 4;
            }
            Instruction::Halt => {
                self.state = State::Halted;
            }
        }
    }

    fn check_state(&self, s: State) -> Result<(), String> {
        if self.state != s {
            return Err(format!("bad machine state: {:?}", self.state));
        }
        Ok(())
    }

    fn step(&mut self) -> Result<(), String> {
        self.check_state(State::Running)?;
        let instr = self.get_instr()?;
        self.exec(&instr);
        Ok(())
    }

    pub fn run(&mut self) -> Result<State, String> {
        while self.state == State::Running {
            self.step()?;
        }
        Ok(self.state)
    }
}
