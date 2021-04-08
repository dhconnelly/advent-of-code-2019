import { readFileSync } from "fs";
import { VM, State } from "../intcode.js";

class AmpCircuit {
    constructor(prog, n, seq) {
        this.signal = 0;
        this.vms = [];
        for (let i = 0; i < n; i++) {
            let vm = new VM(prog);
            vm.run();
            vm.write(seq[i]);
            vm.run();
            this.vms.push(vm);
        }
    }

    step() {
        for (let vm of this.vms) {
            if (vm.state === State.HALT) return;
            vm.write(this.signal);
            vm.run();
            this.signal = vm.read();
            vm.run();
        }
    }

    is_halted() {
        return this.vms[this.vms.length - 1].state === State.HALT;
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

function maximize(prog, phases) {
    let maxSignal = Number.MIN_SAFE_INTEGER;
    genPhaseSeq(phases.length, phases, [], (seq) => {
        let circuit = new AmpCircuit(prog, 5, seq);
        circuit.run();
        maxSignal = Math.max(maxSignal, circuit.signal);
    });
    return maxSignal;
}

export function main(path) {
    const file = readFileSync(path, "ascii");
    const toks = file.split(",");
    const prog = toks.map((s) => parseInt(s, 10));
    console.log(maximize(prog, [0, 1, 2, 3, 4]));
    console.log(maximize(prog, [5, 6, 7, 8, 9]));
}
