const fs = require("fs");

const State = {
    HALT: Symbol("halt"),
    RUN: Symbol("run"),
};

class VM {
    constructor(prog) {
        this.state = State.RUN;
        this.mem = prog.slice();
        this.pc = 0;
    }

    step() {
        let mem = this.mem;
        let pc = this.pc;
        switch (mem[pc]) {
            case 1:
                mem[mem[pc + 3]] = mem[mem[pc + 1]] + mem[mem[pc + 2]];
                this.pc += 4;
                break;
            case 2:
                mem[mem[pc + 3]] = mem[mem[pc + 1]] * mem[mem[pc + 2]];
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

function run(prog, a, b) {
    prog = prog.slice();
    prog[1] = a;
    prog[2] = b;
    let vm = new VM(prog);
    vm.run();
    return vm.mem[0];
}

function main(argv) {
    const path = argv[0];
    const file = fs.readFileSync(path, "ascii");
    const toks = file.split(",");
    const prog = toks.map((s) => parseInt(s, 10));
    console.log(run(prog, 12, 2));
}

main(process.argv.slice(2));
