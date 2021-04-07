const fs = require("fs");

function fuelFor(mass) {
    return Math.floor(mass / 3) - 2;
}

function recFuelFor(mass) {
    let total = 0;
    for (let fuel = fuelFor(mass); fuel > 0; fuel = fuelFor(fuel)) {
        total += fuel;
    }
    return total;
}

function main(argv) {
    const path = argv[0];
    const file = fs.readFileSync(path, "ascii");
    const lines = file.split("\n");
    lines.pop();

    const masses = lines.map((s) => parseInt(s, 10));
    const fuels = masses.map(fuelFor);
    console.log(fuels.reduce((acc, x) => acc + x));

    const recFuels = masses.map(recFuelFor);
    console.log(recFuels.reduce((acc, x) => acc + x));
}

main(process.argv.slice(2));
