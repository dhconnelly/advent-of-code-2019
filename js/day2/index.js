import { readFileSync } from "fs";
import { VM } from "../intcode.js";

function run(prog, a, b) {
    prog = prog.slice();
    prog[1] = a;
    prog[2] = b;
    let vm = new VM(prog);
    vm.run();
    return vm.mem[0];
}

function find(prog, target) {
    for (let noun = 0; noun <= 99; noun++) {
        for (let verb = 0; verb <= 99; verb++) {
            if (run(prog, noun, verb) === target) {
                return 100 * noun + verb;
            }
        }
    }
}

export function main(path) {
    const file = readFileSync(path, "ascii");
    const toks = file.split(",");
    const prog = toks.map((s) => parseInt(s, 10));
    console.log(run(prog, 12, 2));
    console.log(find(prog, 19690720));
}
