"use strict";

const fs = require("fs");
const intcode = require("../intcode");
const util = require("../util");

const Dir = util.makeEnum({
    1: "NORTH",
    2: "SOUTH",
    3: "WEST",
    4: "EAST",
});

let Status = util.makeEnum({
    0: "WALL",
    1: "HALL",
    2: "OXY",
});
Status.toString = function (s) {
    // prettier-ignore
    switch (s) {
        case Status.WALL: return "#";
        case Status.HALL: return ".";
        case Status.OXY: return "@";
        default: return "?";
    }
};

class IntDroid {
    constructor(prog) {
        this.vm = new intcode.VM(prog, { debug: false });
        this.vm.run();
        this.status = null;
    }

    move(dir) {
        util.assertEq(this.vm.state, intcode.State.READ);
        this.vm.write(Dir.int(dir));
        this.vm.run();

        util.assertEq(this.vm.state, intcode.State.WRITE);
        this.status = Status.of(this.vm.read());
        this.vm.run();

        return this.status;
    }
}

function go(pos, dir) {
    // prettier-ignore
    switch (dir) {
        case Dir.NORTH: return { x: pos.x, y: pos.y + 1 };
        case Dir.SOUTH: return { x: pos.x, y: pos.y - 1 };
        case Dir.EAST: return { x: pos.x + 1, y: pos.y };
        case Dir.WEST: return { x: pos.x - 1, y: pos.y };
    }
    util.assert(false);
}

function opposite(dir) {
    // prettier-ignore
    switch (dir) {
        case Dir.NORTH: return Dir.SOUTH;
        case Dir.SOUTH: return Dir.NORTH;
        case Dir.EAST: return Dir.WEST;
        case Dir.WEST: return Dir.EAST;
    }
    util.assert(false);
}

function keyFor(pt) {
    return `${pt.x},${pt.y}`;
}

function fromKey(k) {
    let toks = k.split(",");
    return { x: +toks[0], y: +toks[1] };
}

class Explorer {
    constructor(prog) {
        this.droid = new IntDroid(prog);
        this.map = new Map();
    }

    explore(pos) {
        // starting |pos|, explore in all directions
        for (let dir of Dir.all()) {
            let newPos = go(pos, dir);
            let k = keyFor(newPos);
            if (this.map.has(k)) continue;

            let status = this.droid.move(dir);
            this.map.set(k, status);

            // if we moved into a hallway, explore from there, but move back
            // when finished
            if (status !== Status.WALL) {
                this.explore(newPos);
                let retStatus = this.droid.move(opposite(dir));
                util.assertEq(Status.HALL, retStatus);
            }
        }
    }
}

function explore(prog) {
    let e = new Explorer(prog);
    e.explore({ x: 0, y: 0 });
    return e.map;
}

const INT_MAX = Number.MAX_SAFE_INTEGER;
const INT_MIN = Number.MIN_SAFE_INTEGER;

function printMap(map) {
    // prettier-ignore
    let minX = INT_MAX, minY = INT_MAX, maxX = INT_MIN, maxY = INT_MIN;
    for (let e of map) {
        let pt = fromKey(e[0]);
        if (pt.x < minX) minX = pt.x;
        if (pt.x > maxX) maxX = pt.x;
        if (pt.y < minY) minY = pt.y;
        if (pt.y > maxY) maxY = pt.y;
    }
    for (let y = minY; y <= maxY; y++) {
        let row = "";
        for (let x = minX; x <= maxX; x++) {
            row += Status.toString(map.get(keyFor({ x: x, y: y })));
        }
        console.log(row);
    }
}

function bfs(map, start) {
    let d = new Map();
    d.set(keyFor(start), 0);
    let v = new Set();
    v.add(keyFor(start));
    let q = [{ pos: start, dist: 0 }];
    while (q.length !== 0) {
        let front = q.shift();
        for (let dir of Dir.all()) {
            let nbr = go(front.pos, dir);
            let nbrk = keyFor(nbr);
            if (v.has(nbrk)) continue;
            v.add(nbrk);
            if (map.get(nbrk) !== Status.WALL) d.set(nbrk, front.dist + 1);
            if (map.get(nbrk) !== Status.HALL) continue;
            q.push({ pos: nbr, dist: front.dist + 1 });
        }
    }
    return d;
}

function find(map, target) {
    for (let e of map) if (e[1] === target) return fromKey(e[0]);
}

function max(dists) {
    let max = INT_MIN;
    for (let e of dists) if (e[1] > max) max = e[1];
    return max;
}

function main(args) {
    const path = args[0];
    const file = fs.readFileSync(path, "ascii");
    const toks = file.split(",");
    const prog = toks.map((s) => parseInt(s, 10));
    const map = explore(prog);

    const startDists = bfs(map, { x: 0, y: 0 });
    const oxy = find(map, Status.OXY);
    console.log(startDists.get(keyFor(oxy)));

    const oxyDists = bfs(map, oxy);
    console.log(max(oxyDists));
}

main(process.argv.slice(2));
