"use strict";

function omap(obj, f) {
    return Object.fromEntries(Object.entries(obj).map((e) => f(e[0], e[1])));
}

function makeEnum(enumMap) {
    let intToSymbol = omap(enumMap, (k, v) => [k, Symbol(v)]);
    let symbolToInt = omap(intToSymbol, (k, v) => [v, k]);
    let nameToSymbol = omap(intToSymbol, (k, v) => [enumMap[k], v]);
    let symbolToName = omap(nameToSymbol, (k, v) => [v, k]);
    nameToSymbol.of = (i) => intToSymbol[i];
    nameToSymbol.str = (sym) => symbolToName[sym];
    nameToSymbol.int = (sym) => symbolToInt[sym];
    return nameToSymbol;
}

exports.makeEnum = makeEnum;
exports.omap = omap;
