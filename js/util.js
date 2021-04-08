"use strict";

function omap(obj, f) {
    return Object.fromEntries(Object.entries(obj).map((e) => f(e[0], e[1])));
}

function makeEnum(enumMap) {
    let intToSymbol = omap(enumMap, (k, v) => [k, Symbol(v)]);
    let symbolToInt = omap(intToSymbol, (k, v) => [v, +k]);
    let nameToSymbol = omap(intToSymbol, (k, v) => [enumMap[k], v]);
    let symbolToName = omap(nameToSymbol, (k, v) => [v, k]);
    let eenum = {};
    Object.assign(eenum, nameToSymbol);
    eenum.of = (i) => intToSymbol[i];
    eenum.str = (sym) => symbolToName[sym];
    eenum.int = (sym) => symbolToInt[sym];
    eenum.all = () => Object.getOwnPropertySymbols(symbolToName);
    return eenum;
}

function assert(cond) {
    if (!cond) {
        throw new Error("assertion failed:", cond);
    }
}

function assertEq(a, b) {
    if (a !== b) {
        throw new Error("assertion failed:", a, "===", b);
    }
}

exports.makeEnum = makeEnum;
exports.omap = omap;
exports.assert = assert;
exports.assertEq = assert;
