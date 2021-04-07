"use strict";

const fs = require("fs");
const intcode = require("../intcode");

function run(prog, input) {
    const vm = new intcode.VM(
        prog,
        () => input,
        (x) => console.log(x)
    );
    vm.run();
}

function main(args) {
    const path = args[0];
    const file = fs.readFileSync(path, "ascii");
    const toks = file.split(",");
    const prog = toks.map((s) => parseInt(s, 10));
    run(prog, 1);
    run(prog, 2);
}

main(process.argv.slice(2));
