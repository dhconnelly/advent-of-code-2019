const days = ["day1", "day2", "day5", "day7", "day9", "day15"];

const loadModule = async (day) => [day, await import(`./${day}/index.js`)];
const loadModules = async () =>
    await Promise.all(days.map((day) => loadModule(day)));
const modules = Object.fromEntries(await loadModules());

function run(day) {
    console.log(">", day);
    modules[day].main(`./${day}/input.txt`);
    console.log();
}

function runAll() {
    for (let day of days) run(day);
}

function main(args) {
    if (args.length === 0) runAll();
    else run(args[0]);
}

main(process.argv.slice(2));
