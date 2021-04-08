"use strict";

function omap(obj, f) {
    return Object.fromEntries(Object.entries(obj).map((e) => f(e[0], e[1])));
}

function makeEnum(enumMap) {
    let intToSymbol = omap(enumMap, (k, v) => [k, Symbol(v)]);
    let nameToSymbol = omap(intToSymbol, (k, v) => [enumMap[k], v]);
    let symbolToName = omap(nameToSymbol, (k, v) => [v, k]);
    nameToSymbol.of = (i) => intToSymbol[i];
    nameToSymbol.str = (sym) => symbolToName[sym];
    return nameToSymbol;
}

const State = makeEnum({
    0: "HALT",
    1: "RUN",
    2: "READ",
    3: "WRITE",
});

const Mode = makeEnum({
    0: "POS",
    1: "IMM",
    2: "REL",
});

const Opcode = makeEnum({
    1: "ADD",
    2: "MUL",
    3: "READ",
    4: "WRITE",
    5: "JMPIF",
    6: "JMPNOT",
    7: "LT",
    8: "EQ",
    9: "ADJREL",
    99: "HALT",
});

class ExecutionError extends Error {
    constructor(msg) {
        super(msg);
        this.name = "ExecutionError";
    }
}

function div(x, y) {
    return Math.floor(x / y);
}

class VM {
    constructor(prog, opts) {
        this.mem = prog.slice();
        this.debug = opts && !!opts.debug;
        this.state = State.RUN;
        this.input = 0;
        this.output = 0;
        this.sp = 0;
        this.pc = 0;
    }

    error(msg) {
        throw new ExecutionError(`execution error at pc=${this.pc}: ${msg}`);
    }

    nextOp() {
        let op = this.mem[this.pc];
        let code = Opcode.of(op % 100);
        let modes = div(op, 100);
        let mode1 = Mode.of(modes % 10);
        let mode2 = Mode.of(div(modes, 10) % 10);
        let mode3 = Mode.of(div(modes, 100) % 10);
        return {
            code: code,
            modes: [mode1, mode2, mode3],
        };
    }

    get(arg, mode) {
        const base = this.pc + 1;
        switch (mode) {
            case Mode.IMM:
                return this.mem[base + arg];
            case Mode.POS:
                return this.mem[this.mem[base + arg]];
            case Mode.REL:
                return this.mem[this.sp + this.mem[base + arg]];
        }
    }

    set(arg, mode, val) {
        const base = this.pc + 1;
        switch (mode) {
            case Mode.POS:
                this.mem[this.mem[base + arg]] = val;
                break;
            case Mode.IMM:
                this.error("can't write in immediate mode");
                break;
            case Mode.REL:
                this.mem[this.sp + this.mem[base + arg]] = val;
                break;
        }
    }

    write(x) {
        let op = this.nextOp();
        this.set(0, op.modes[0], x);
        this.state = State.RUN;
        this.pc += 2;
    }

    read() {
        this.state = State.RUN;
        return this.output;
    }

    step() {
        let op = this.nextOp();
        if (this.debug) {
            console.log(`pc=${this.pc}\t${Opcode.str(op.code)}`);
            console.log(this.mem);
        }
        let modes = op.modes;
        let a = this.get(0, modes[0]);
        let b = this.get(1, modes[1]);

        switch (op.code) {
            case Opcode.ADD:
                this.set(2, modes[2], a + b);
                this.pc += 4;
                break;

            case Opcode.MUL:
                this.set(2, modes[2], a * b);
                this.pc += 4;
                break;

            case Opcode.READ:
                this.state = State.READ;
                break;

            case Opcode.WRITE:
                this.output = a;
                this.state = State.WRITE;
                this.pc += 2;
                break;

            case Opcode.JMPIF:
                if (a !== 0) this.pc = b;
                else this.pc += 3;
                break;

            case Opcode.JMPNOT:
                if (a === 0) this.pc = b;
                else this.pc += 3;
                break;

            case Opcode.LT:
                this.set(2, modes[2], a < b ? 1 : 0);
                this.pc += 4;
                break;

            case Opcode.EQ:
                this.set(2, modes[2], a === b ? 1 : 0);
                this.pc += 4;
                break;

            case Opcode.ADJREL:
                this.sp += a;
                this.pc += 2;
                break;

            case Opcode.HALT:
                this.state = State.HALT;
                break;
        }
    }

    run() {
        while (this.state === State.RUN) this.step();
    }
}

exports.State = State;
exports.VM = VM;
