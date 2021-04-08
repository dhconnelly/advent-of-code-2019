import { readFileSync } from "fs";
import { VM, State } from "../intcode.js";

function run(prog, testId) {
    let diagCode;
    const vm = new VM(prog);
    while (vm.state !== State.HALT) {
        switch (vm.state) {
            case State.RUN:
                vm.run();
                break;
            case State.WRITE:
                diagCode = vm.read();
                break;
            case State.READ:
                vm.write(testId);
                break;
        }
    }
    return diagCode;
}

export function main(path) {
    const file = readFileSync(path, "ascii");
    const toks = file.split(",");
    const prog = toks.map((s) => parseInt(s, 10));
    console.log(run(prog, 1));
    console.log(run(prog, 5));
}
