"use strict";

const fs = require("fs");
const intcode = require("../intcode");

class AmpCircuit {
    constructor(prog, n, seq) {
        this.signal = 0;
        this.vms = [];
        for (let i = 0; i < n; i++) {
            let vm = new intcode.VM(prog);
            vm.run();
            vm.write(seq[i]);
            vm.run();
            this.vms.push(vm);
        }
    }

    step() {
        for (let vm of this.vms) {
            if (vm.state === intcode.State.HALT) return;
            vm.write(this.signal);
            vm.run();
            this.signal = vm.read();
            vm.run();
        }
    }

    is_halted() {
        return this.vms[this.vms.length - 1].state === intcode.State.HALT;
    }

    run() {
        while (!this.is_halted()) this.step();
    }
}

function genPhaseSeq(n, phases, cur, f) {
    if (cur.length === n) return f(cur);
    for (let i = 0; i < phases.length; i++) {
        let phase = phases[i];
        if (phase === undefined) continue;
        phases[i] = undefined;
        genPhaseSeq(n, phases, cur.concat(phase), f);
        phases[i] = phase;
    }
}

function main(args) {
    const path = args[0];
    const file = fs.readFileSync(path, "ascii");
    const toks = file.split(",");
    const prog = toks.map((s) => parseInt(s, 10));
    let maxSignal = Number.MIN_SAFE_INTEGER;
    genPhaseSeq(5, [0, 1, 2, 3, 4], [], (seq) => {
        let circuit = new AmpCircuit(prog, 5, seq);
        circuit.run();
        maxSignal = Math.max(maxSignal, circuit.signal);
    });
    console.log(maxSignal);
}

main(process.argv.slice(2));
