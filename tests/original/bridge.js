const { spawnSync } = require('child_process');
const path = require('path');

const CLI_PATH = path.join(__dirname, '../../bin/decimal-cli');

function callGo(op, args) {
  const payload = JSON.stringify({
    op,
    precision: Decimal.precision,
    rounding: Decimal.rounding,
    toExpNeg: Decimal.toExpNeg,
    toExpPos: Decimal.toExpPos,
    minE: Decimal.minE,
    maxE: Decimal.maxE,
    args
  });

  const res = spawnSync(CLI_PATH, ['--rpc'], {
    input: payload + '\n',
    encoding: 'utf8'
  });

  if (res.error) {
    throw new Error('Go CLI Execution Error: ' + res.error.message);
  }

  const output = res.stdout.trim();
  if (!output) {
    throw new Error('Go CLI returned empty output');
  }

  return JSON.parse(output);
}

function Decimal(val) {
  if (!(this instanceof Decimal)) {
    return new Decimal(val);
  }

  const res = callGo('new', [val]);
  this.s = res.s;
  this.e = res.e;
  this.d = res.d;
  this._str = res.str;
}

Decimal.precision = 20;
Decimal.rounding = 4;
Decimal.toExpNeg = -7;
Decimal.toExpPos = 21;
Decimal.minE = -9e15;
Decimal.maxE = 9e15;

Decimal.config = function (obj) {
  if (obj) {
    if (obj.precision !== undefined) Decimal.precision = obj.precision;
    if (obj.rounding !== undefined) Decimal.rounding = obj.rounding;
    if (obj.toExpNeg !== undefined) Decimal.toExpNeg = obj.toExpNeg;
    if (obj.toExpPos !== undefined) Decimal.toExpPos = obj.toExpPos;
    if (obj.minE !== undefined) Decimal.minE = obj.minE;
    if (obj.maxE !== undefined) Decimal.maxE = obj.maxE;
  }
  return Decimal;
};

Decimal.set = Decimal.config;

function wrapRes(res) {
  const d = new Decimal(0);
  d.s = res.s;
  d.e = res.e;
  d.d = res.d;
  d._str = res.str;
  return d;
}

Decimal.prototype.plus = Decimal.prototype.add = function (y) {
  return wrapRes(callGo('plus', [this.valueOf(), y instanceof Decimal ? y.valueOf() : y]));
};

Decimal.prototype.minus = Decimal.prototype.sub = function (y) {
  return wrapRes(callGo('minus', [this.valueOf(), y instanceof Decimal ? y.valueOf() : y]));
};

Decimal.prototype.times = Decimal.prototype.mul = function (y) {
  return wrapRes(callGo('times', [this.valueOf(), y instanceof Decimal ? y.valueOf() : y]));
};

Decimal.prototype.dividedBy = Decimal.prototype.div = function (y) {
  return wrapRes(callGo('dividedBy', [this.valueOf(), y instanceof Decimal ? y.valueOf() : y]));
};

Decimal.prototype.modulo = Decimal.prototype.mod = function (y) {
  return wrapRes(callGo('mod', [this.valueOf(), y instanceof Decimal ? y.valueOf() : y]));
};

Decimal.prototype.squareRoot = Decimal.prototype.sqrt = function () {
  return wrapRes(callGo('sqrt', [this.valueOf()]));
};

Decimal.prototype.cubeRoot = Decimal.prototype.cbrt = function () {
  return wrapRes(callGo('cbrt', [this.valueOf()]));
};

Decimal.prototype.toPower = Decimal.prototype.pow = function (y) {
  return wrapRes(callGo('pow', [this.valueOf(), y instanceof Decimal ? y.valueOf() : y]));
};

Decimal.prototype.naturalLogarithm = Decimal.prototype.ln = function () {
  return wrapRes(callGo('ln', [this.valueOf()]));
};

Decimal.prototype.naturalExponential = Decimal.prototype.exp = function () {
  return wrapRes(callGo('exp', [this.valueOf()]));
};

Decimal.prototype.sine = Decimal.prototype.sin = function () {
  return wrapRes(callGo('sin', [this.valueOf()]));
};

Decimal.prototype.cosine = Decimal.prototype.cos = function () {
  return wrapRes(callGo('cos', [this.valueOf()]));
};

Decimal.prototype.tangent = Decimal.prototype.tan = function () {
  return wrapRes(callGo('tan', [this.valueOf()]));
};

Decimal.prototype.abs = Decimal.prototype.absoluteValue = function () {
  return wrapRes(callGo('abs', [this.valueOf()]));
};

Decimal.prototype.neg = Decimal.prototype.negated = function () {
  return wrapRes(callGo('neg', [this.valueOf()]));
};

Decimal.prototype.trunc = Decimal.prototype.truncated = function () {
  return wrapRes(callGo('sub', [this.valueOf(), '0']));
};

Decimal.prototype.cmp = Decimal.prototype.comparedTo = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('cmp', [this.valueOf(), yVal]);
  return res.isCmp ? res.cmpRes : NaN;
};

Decimal.prototype.eq = Decimal.prototype.equals = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('eq', [this.valueOf(), yVal]);
  return res.isBool ? res.boolRes : false;
};

Decimal.prototype.gt = Decimal.prototype.greaterThan = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('gt', [this.valueOf(), yVal]);
  return res.isBool ? res.boolRes : false;
};

Decimal.prototype.gte = Decimal.prototype.greaterThanOrEqualTo = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('gte', [this.valueOf(), yVal]);
  return res.isBool ? res.boolRes : false;
};

Decimal.prototype.lt = Decimal.prototype.lessThan = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('lt', [this.valueOf(), yVal]);
  return res.isBool ? res.boolRes : false;
};

Decimal.prototype.lte = Decimal.prototype.lessThanOrEqualTo = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('lte', [this.valueOf(), yVal]);
  return res.isBool ? res.boolRes : false;
};

Decimal.prototype.isNaN = function () {
  return this.s === 0 || this.s === null || (this.d === null && this.s === 0);
};

Decimal.prototype.isZero = function () {
  return this.d && this.d.length > 0 && this.d[0] === 0;
};

Decimal.prototype.isFinite = function () {
  return this.d !== null && this.s !== 0;
};

Decimal.prototype.toString = Decimal.prototype.valueOf = function () {
  return this._str || '0';
};

module.exports = Decimal;
