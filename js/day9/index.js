import { readFileSync } from "fs";
import { VM, State } from "../intcode.js";

function run(prog, input) {
    const vm = new VM(prog);
    while (vm.state !== State.HALT) {
        switch (vm.state) {
            case State.READ:
                vm.write(input);
                break;
            case State.WRITE:
                console.log(vm.read());
                break;
            case State.RUN:
                vm.run();
                break;
        }
    }
}

export function main(path) {
    const file = readFileSync(path, "ascii");
    const toks = file.split(",");
    const prog = toks.map((s) => parseInt(s, 10));
    run(prog, 1);
    run(prog, 2);
}
