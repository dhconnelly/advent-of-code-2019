"use strict";

const fs = require("fs");
const intcode = require("../intcode");

function run(prog, input) {
    const vm = new intcode.VM(prog);
    while (vm.state !== intcode.State.HALT) {
        switch (vm.state) {
            case intcode.State.READ:
                vm.write(input);
                break;
            case intcode.State.WRITE:
                console.log(vm.read());
                break;
            case intcode.State.RUN:
                vm.run();
                break;
        }
    }
}

function main(path) {
    const file = fs.readFileSync(path, "ascii");
    const toks = file.split(",");
    const prog = toks.map((s) => parseInt(s, 10));
    run(prog, 1);
    run(prog, 2);
}

module.exports = main;
