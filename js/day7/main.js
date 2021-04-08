"use strict";

const fs = require("fs");
const intcode = require("../intcode");

function run(prog, input1, input2) {
    let vm = new intcode.VM(prog);
    vm.run();
    vm.write(input1);
    vm.run();
    vm.write(input2);
    vm.run();
    return vm.read();
}

function runCircuit(prog, n, seq) {
    let signal = 0;
    for (let i = 0; i < n; i++) {
        signal = run(prog, seq[i], signal);
    }
    return signal;
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
        let signal = runCircuit(prog, 5, seq);
        maxSignal = Math.max(maxSignal, signal);
    });
    console.log(maxSignal);
}

main(process.argv.slice(2));
