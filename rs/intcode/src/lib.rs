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
    Read,
    Write,
    JmpIf,
    JmpNot,
    Lt,
    Eq,
    AdjRel,
    Halt,
}

impl Opcode {
    fn new(x: i64) -> Result<Opcode, String> {
        match x {
            1 => Ok(Opcode::Add),
            2 => Ok(Opcode::Mul),
            3 => Ok(Opcode::Read),
            4 => Ok(Opcode::Write),
            5 => Ok(Opcode::JmpIf),
            6 => Ok(Opcode::JmpNot),
            7 => Ok(Opcode::Lt),
            8 => Ok(Opcode::Eq),
            9 => Ok(Opcode::AdjRel),
            99 => Ok(Opcode::Halt),
            _ => Err(format!("unknown opcode: {}", x)),
        }
    }
}

#[derive(Clone, Copy, Debug)]
enum Arg {
    Pos(i64),
    Imm(i64),
    Rel(i64),
}

impl Arg {
    fn new(x: i64, mode: i64) -> Result<Arg, String> {
        match mode {
            0 => Ok(Arg::Pos(x)),
            1 => Ok(Arg::Imm(x)),
            2 => Ok(Arg::Rel(x)),
            _ => Err(format!("bad param mode: {}", mode)),
        }
    }

    fn args(x: i64, y: i64, z: i64, modes: i64) -> Result<(Arg, Arg, Arg), String> {
        Ok((
            Arg::new(x, modes % 10)?,
            Arg::new(y, (modes / 10) % 10)?,
            Arg::new(z, (modes / 100) % 10)?,
        ))
    }
}

#[derive(Debug, Copy, Clone)]
enum Instruction {
    Add(Arg, Arg, Arg),
    Mul(Arg, Arg, Arg),
    Read(Arg),
    Write(Arg),
    JmpIf(Arg, Arg),
    JmpNot(Arg, Arg),
    Lt(Arg, Arg, Arg),
    Eq(Arg, Arg, Arg),
    AdjRel(Arg),
    Halt,
}

#[derive(Debug, PartialEq, Copy, Clone)]
pub enum State {
    Running,
    Reading,
    Writing,
    Halted,
}

pub struct Machine {
    state: State,
    mem: HashMap<i64, i64>,
    pc: i64,
    rel: i64,
    instr: Option<Instruction>,
    input: i64,
    output: i64,
}

impl Machine {
    pub fn new(program: &Program) -> Machine {
        Machine {
            state: State::Running,
            mem: program.to_memory(),
            instr: None,
            pc: 0,
            rel: 0,
            input: 0,
            output: 0,
        }
    }

    pub fn write(&mut self, x: i64) {
        self.input = x;
    }

    pub fn read(&self) -> i64 {
        self.output
    }

    pub fn get(&self, i: i64) -> i64 {
        *self.mem.get(&i).unwrap_or(&0)
    }

    fn load(&self, arg: &Arg) -> i64 {
        match arg {
            Arg::Pos(i) => self.get(*i),
            Arg::Imm(i) => *i,
            Arg::Rel(i) => self.get(self.rel + *i),
        }
    }

    pub fn set(&mut self, i: i64, x: i64) {
        self.mem.insert(i, x);
    }

    fn store(&mut self, arg: &Arg, x: i64) -> Result<(), String> {
        match arg {
            Arg::Pos(i) => Ok(self.set(*i, x)),
            Arg::Imm(_) => Err(format!("bad mode for write: {:?}", arg)),
            Arg::Rel(i) => Ok(self.set(self.rel + *i, x)),
        }
    }

    fn get_instr(&self) -> Result<Instruction, String> {
        let arg0 = self.get(self.pc + 0);
        let arg1 = self.get(self.pc + 1);
        let arg2 = self.get(self.pc + 2);
        let arg3 = self.get(self.pc + 3);
        let args = Arg::args(arg1, arg2, arg3, arg0 / 100)?;
        let instr = match Opcode::new(arg0 % 100)? {
            Opcode::Add => Instruction::Add(args.0, args.1, args.2),
            Opcode::Mul => Instruction::Mul(args.0, args.1, args.2),
            Opcode::Read => Instruction::Read(args.0),
            Opcode::Write => Instruction::Write(args.0),
            Opcode::JmpIf => Instruction::JmpIf(args.0, args.1),
            Opcode::JmpNot => Instruction::JmpNot(args.0, args.1),
            Opcode::Lt => Instruction::Lt(args.0, args.1, args.2),
            Opcode::Eq => Instruction::Eq(args.0, args.1, args.2),
            Opcode::AdjRel => Instruction::AdjRel(args.0),
            Opcode::Halt => Instruction::Halt,
        };
        Ok(instr)
    }

    fn exec(&mut self, instr: &Instruction) -> Result<(), String> {
        match instr {
            Instruction::Add(x, y, z) => {
                self.store(z, self.load(x) + self.load(y))?;
                self.state = State::Running;
                self.pc += 4;
            }

            Instruction::Mul(x, y, z) => {
                self.store(z, self.load(x) * self.load(y))?;
                self.state = State::Running;
                self.pc += 4;
            }

            Instruction::Read(_) => {
                self.state = State::Reading;
            }

            Instruction::Write(x) => {
                self.output = self.load(x);
                self.state = State::Writing;
                self.pc += 2;
            }

            Instruction::JmpIf(x, y) => {
                if self.load(x) != 0 {
                    self.pc = self.load(y);
                } else {
                    self.pc += 3;
                }
                self.state = State::Running;
            }

            Instruction::JmpNot(x, y) => {
                if self.load(x) == 0 {
                    self.pc = self.load(y);
                } else {
                    self.pc += 3;
                }
                self.state = State::Running;
            }

            Instruction::Lt(x, y, z) => {
                if self.load(x) < self.load(y) {
                    self.store(z, 1)?;
                } else {
                    self.store(z, 0)?;
                }
                self.pc += 4;
                self.state = State::Running;
            }

            Instruction::Eq(x, y, z) => {
                if self.load(x) == self.load(y) {
                    self.store(z, 1)?;
                } else {
                    self.store(z, 0)?;
                }
                self.pc += 4;
                self.state = State::Running;
            }

            Instruction::AdjRel(x) => {
                self.rel += self.load(x);
                self.pc += 2;
                self.state = State::Running;
            }

            Instruction::Halt => {
                self.state = State::Halted;
            }
        }
        Ok(())
    }

    fn step(&mut self) -> Result<(), String> {
        if self.state == State::Halted {
            return Err(format!("bad machine state: {:?}", self.state));
        }
        if self.state == State::Reading {
            let instr = self.instr.ok_or("no instruction while reading")?;
            if let Instruction::Read(x) = instr {
                self.store(&x, self.input)?;
                self.pc += 2;
            } else {
                return Err(format!("non-read instruction while reading: {:?}", instr));
            }
        }
        self.instr = Some(self.get_instr()?);
        self.exec(&self.instr.unwrap())
    }

    pub fn run(&mut self) -> Result<State, String> {
        self.step()?;
        while self.state == State::Running {
            self.step()?;
        }
        Ok(self.state)
    }
}
