const days = ["day1", "day2", "day5", "day7", "day9", "day15"];

function run(day) {
    let base = "./" + day;
    let m = require(base);
    let path = base + "/input.txt";
    console.log(">", day);
    m(path);
    console.log();
}

function runAll() {
    days.forEach(run);
}

function main(args) {
    if (args.length === 0) runAll();
    else run(args[0]);
}

main(process.argv.slice(2));
