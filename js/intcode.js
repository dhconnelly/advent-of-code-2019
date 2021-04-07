"use strict";

const State = Object.freeze({
    HALT: Symbol("halt"),
    RUN: Symbol("run"),
});

const Mode = Object.freeze({
    IMM: Symbol("immediate"),
    POS: Symbol("position"),
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

function toMode(x) {
    switch (x) {
        case 0:
            return Mode.POS;
        case 1:
            return Mode.IMM;
    }
}

class VM {
    constructor(prog, getInput, writeOutput) {
        this.state = State.RUN;
        this.mem = prog.slice();
        this.getInput = getInput;
        this.writeOutput = writeOutput;
        this.pc = 0;
    }

    error(msg) {
        throw new ExecutionError(`execution error at pc=${this.pc}: ${msg}`);
    }

    nextOp() {
        let op = this.mem[this.pc];
        let code = op % 100;
        op = div(op, 100);
        let mode1 = toMode(op % 10);
        op = div(op, 10);
        let mode2 = toMode(op % 10);
        op = div(op, 10);
        let mode3 = toMode(op % 10);
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
        }
    }

    step() {
        let op = this.nextOp();
        let modes = op.modes;
        let a = this.get(0, modes[0]);
        let b = this.get(1, modes[1]);
        switch (op.code) {
            case 1:
                this.set(2, modes[2], a + b);
                this.pc += 4;
                break;

            case 2:
                this.set(2, modes[2], a * b);
                this.pc += 4;
                break;

            case 3:
                this.set(0, modes[0], this.getInput());
                this.pc += 2;
                break;

            case 4:
                this.writeOutput(a);
                this.pc += 2;
                break;

            case 5:
                if (a !== 0) this.pc = b;
                else this.pc += 3;
                break;

            case 6:
                if (a === 0) this.pc = b;
                else this.pc += 3;
                break;

            case 7:
                this.set(2, modes[2], a < b ? 1 : 0);
                this.pc += 4;
                break;

            case 8:
                this.set(2, modes[2], a === b ? 1 : 0);
                this.pc += 4;
                break;

            case 99:
                this.state = State.HALT;
                break;
        }
    }

    run() {
        while (this.state != State.HALT) this.step();
    }
}

exports.State = State;
exports.VM = VM;
