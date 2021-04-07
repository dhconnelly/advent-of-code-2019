"use strict";

const fs = require("fs");
const intcode = require("../intcode");

function run(prog, testId) {
    let diagCode;
    const vm = new intcode.VM(
        prog,
        () => testId,
        (x) => (diagCode = x)
    );
    vm.run();
    return diagCode;
}

function main(args) {
    const path = args[0];
    const file = fs.readFileSync(path, "ascii");
    const toks = file.split(",");
    const prog = toks.map((s) => parseInt(s, 10));
    console.log(run(prog, 1));
    console.log(run(prog, 5));
}

main(process.argv.slice(2));
