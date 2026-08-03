const { spawnSync } = require('child_process');
const path = require('path');
const DecimalJS = require('../../../decimal.js/decimal.js');

const CLI_PATH = path.join(__dirname, '../../bin/decimal-cli');

function callGo(op, args, Ctor = Decimal) {
  const serializedArgs = (args || []).map(a => {
    if (a === undefined) return 'NaN';
    if (typeof a === 'number' && isNaN(a)) return 'NaN';
    if (a === Infinity) return 'Infinity';
    if (a === -Infinity) return '-Infinity';
    if (Object.is(a, -0)) return '-0';
    return String(a);
  });

  const payload = JSON.stringify({
    op,
    precision: Ctor.precision,
    rounding: Ctor.rounding,
    modulo: Ctor.modulo,
    toExpNeg: Ctor.toExpNeg,
    toExpPos: Ctor.toExpPos,
    minE: Ctor.minE,
    maxE: Ctor.maxE,
    args: serializedArgs
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

function callGoCtor(Ctor, op, args) {
  return callGo(op, args, Ctor);
}

function Decimal(val) {
  if (!(this instanceof Decimal)) {
    return new Decimal(val);
  }
  if (val === undefined || typeof val === 'function') {
    throw new Error('[decimal.js] DecimalError: Invalid argument: ' + val);
  }

  const res = callGo('new', [val], this.constructor);
  if (res.error) {
    throw new Error('[decimal.js] DecimalError: ' + res.error);
  }
  this.s = res.s;
  this.e = res.e;
  this.d = res.d;
  this._str = res.str;
}

Decimal.precision = 20;
Decimal.rounding = 4;
Decimal.modulo = 1;
Decimal.toExpNeg = -7;
Decimal.toExpPos = 21;
Decimal.minE = -9e15;
Decimal.maxE = 9e15;
Decimal.config = Decimal.set = function (obj) {
  if (obj === undefined || obj === null || typeof obj !== 'object') {
    throw new Error('[decimal.js] DecimalError: Invalid argument: ' + obj);
  }

  if (obj.defaults === true) {
    Decimal.precision = 20;
    Decimal.rounding = 4;
    Decimal.toExpNeg = -7;
    Decimal.toExpPos = 21;
    Decimal.minE = -9e15;
    Decimal.maxE = 9e15;
    Decimal.modulo = 1;
    Decimal.crypto = false;
  }

  if (obj.precision !== undefined) {
    if (typeof obj.precision !== 'number' || obj.precision !== Math.floor(obj.precision) || obj.precision < 1 || obj.precision > 1e9) {
      throw new Error('[decimal.js] DecimalError: Invalid precision: ' + obj.precision);
    }
    Decimal.precision = obj.precision;
  }

  if (obj.rounding !== undefined) {
    if (typeof obj.rounding !== 'number' || obj.rounding !== Math.floor(obj.rounding) || obj.rounding < 0 || obj.rounding > 8) {
      throw new Error('[decimal.js] DecimalError: Invalid rounding: ' + obj.rounding);
    }
    Decimal.rounding = obj.rounding;
  }

  if (obj.modulo !== undefined) {
    if (typeof obj.modulo !== 'number' || obj.modulo !== Math.floor(obj.modulo) || obj.modulo < 0 || obj.modulo > 9) {
      throw new Error('[decimal.js] DecimalError: Invalid modulo: ' + obj.modulo);
    }
    Decimal.modulo = obj.modulo;
  }

  if (obj.toExpNeg !== undefined) {
    if (typeof obj.toExpNeg !== 'number' || obj.toExpNeg !== Math.floor(obj.toExpNeg) || obj.toExpNeg > 0 || obj.toExpNeg < -9e15) {
      throw new Error('[decimal.js] DecimalError: Invalid toExpNeg: ' + obj.toExpNeg);
    }
    Decimal.toExpNeg = obj.toExpNeg;
  }

  if (obj.toExpPos !== undefined) {
    if (typeof obj.toExpPos !== 'number' || obj.toExpPos !== Math.floor(obj.toExpPos) || obj.toExpPos < 0 || obj.toExpPos > 9e15) {
      throw new Error('[decimal.js] DecimalError: Invalid toExpPos: ' + obj.toExpPos);
    }
    Decimal.toExpPos = obj.toExpPos;
  }

  if (obj.minE !== undefined) {
    if (typeof obj.minE !== 'number' || obj.minE !== Math.floor(obj.minE) || obj.minE > 0 || obj.minE < -9e15) {
      throw new Error('[decimal.js] DecimalError: Invalid minE: ' + obj.minE);
    }
    Decimal.minE = obj.minE;
  }

  if (obj.maxE !== undefined) {
    if (typeof obj.maxE !== 'number' || obj.maxE !== Math.floor(obj.maxE) || obj.maxE < 0 || obj.maxE > 9e15) {
      throw new Error('[decimal.js] DecimalError: Invalid maxE: ' + obj.maxE);
    }
    Decimal.maxE = obj.maxE;
  }

  if (obj.crypto !== undefined) {
    if (obj.crypto !== true && obj.crypto !== false && obj.crypto !== 1 && obj.crypto !== 0) {
      throw new Error('[decimal.js] DecimalError: Invalid crypto: ' + obj.crypto);
    }
    Decimal.crypto = !!obj.crypto;
  }

  return this || Decimal;
};

Decimal.set = Decimal.config;

const MAX_DIGITS = 1e9;

function checkInt32(n, min, max) {
  if (typeof n !== 'number' || n !== Math.floor(n) || n < min || n > max) {
    throw new Error('[decimal.js] DecimalError: Invalid argument: ' + n);
  }
}

function wrapRes(res) {
  const d = Object.create(Decimal.prototype);
  d.s = res.s;
  d.e = res.e;
  d.d = res.d;
  d._str = res.str;
  return d;
}

function toJS(x) {
  // Sync configuration properties to DecimalJS
  DecimalJS.config({
    precision: x.constructor.precision,
    rounding: x.constructor.rounding,
    toExpNeg: x.constructor.toExpNeg,
    toExpPos: x.constructor.toExpPos,
    minE: x.constructor.minE,
    maxE: x.constructor.maxE
  });
  const d = new DecimalJS(0);
  d.s = x.s;
  d.e = x.e;
  d.d = x.d ? x.d.slice() : null;
  return d;
}

function wrapJS(x, Ctor) {
  if (x === undefined) return new Ctor(NaN);
  if (typeof x !== 'object' || x === null) {
    return new Ctor(x);
  }
  if (!x.d) {
    return wrapRes({ s: x.s, e: 0, d: null, str: x.toString() });
  }
  if (x.d.length === 1 && x.d[0] === 0 && x.s === -1) {
    return wrapRes({ s: -1, e: 0, d: [0], str: "-0" });
  }
  return new Ctor(x.toString());
}

// Decimal.clone implementation
Decimal.clone = function (configObj) {
  if (configObj !== undefined && (configObj === null || typeof configObj !== 'object')) {
    throw new Error('[decimal.js] DecimalError: Invalid argument: ' + configObj);
  }
  const Parent = this || Decimal;
  const Ctor = function (val) {
    if (!(this instanceof Ctor)) {
      return new Ctor(val);
    }
    if (val === undefined || typeof val === 'function') {
      throw new Error('[decimal.js] DecimalError: Invalid argument: ' + val);
    }
    const res = callGoCtor(Ctor, 'new', [val]);
    if (res.error) {
      throw new Error('[decimal.js] DecimalError: ' + res.error);
    }
    this.s = res.s;
    this.e = res.e;
    this.d = res.d;
    this._str = res.str;
    this.constructor = Ctor;
  };

  Ctor.prototype = Decimal.prototype;
  Object.setPrototypeOf(Ctor, Parent);

  const defaults = {
    precision: 20,
    rounding: 4,
    modulo: 1,
    toExpNeg: -7,
    toExpPos: 21,
    minE: -9e15,
    maxE: 9e15
  };

  const isDefaults = configObj && configObj.defaults === true;
  const ps = ['precision', 'rounding', 'modulo', 'toExpNeg', 'toExpPos', 'minE', 'maxE'];
  for (let p of ps) {
    Ctor[p] = (isDefaults ? defaults[p] : Parent[p]);
  }

  Ctor.config = Ctor.set = function (obj) {
    if (obj != null) {
      if (typeof obj !== 'object') {
        throw new Error('[decimal.js] DecimalError: Invalid argument: ' + obj);
      }
      if (obj.precision !== undefined) {
        if (typeof obj.precision !== 'number' || obj.precision !== Math.floor(obj.precision) || obj.precision < 1 || obj.precision > 1e9) {
          throw new Error('[decimal.js] DecimalError: Invalid precision: ' + obj.precision);
        }
        Ctor.precision = obj.precision;
      }
      if (obj.rounding !== undefined) {
        if (typeof obj.rounding !== 'number' || obj.rounding !== Math.floor(obj.rounding) || obj.rounding < 0 || obj.rounding > 8) {
          throw new Error('[decimal.js] DecimalError: Invalid rounding: ' + obj.rounding);
        }
        Ctor.rounding = obj.rounding;
      }
      if (obj.modulo !== undefined) {
        if (typeof obj.modulo !== 'number' || obj.modulo !== Math.floor(obj.modulo) || obj.modulo < 0 || obj.modulo > 9) {
          throw new Error('[decimal.js] DecimalError: Invalid modulo: ' + obj.modulo);
        }
        Ctor.modulo = obj.modulo;
      }
      if (obj.toExpNeg !== undefined) {
        if (typeof obj.toExpNeg !== 'number' || obj.toExpNeg !== Math.floor(obj.toExpNeg) || obj.toExpNeg > 0 || obj.toExpNeg < -9e15) {
          throw new Error('[decimal.js] DecimalError: Invalid toExpNeg: ' + obj.toExpNeg);
        }
        Ctor.toExpNeg = obj.toExpNeg;
      }
      if (obj.toExpPos !== undefined) {
        if (typeof obj.toExpPos !== 'number' || obj.toExpPos !== Math.floor(obj.toExpPos) || obj.toExpPos < 0 || obj.toExpPos > 9e15) {
          throw new Error('[decimal.js] DecimalError: Invalid toExpPos: ' + obj.toExpPos);
        }
        Ctor.toExpPos = obj.toExpPos;
      }
      if (obj.minE !== undefined) {
        if (typeof obj.minE !== 'number' || obj.minE !== Math.floor(obj.minE) || obj.minE > 0 || obj.minE < -9e15) {
          throw new Error('[decimal.js] DecimalError: Invalid minE: ' + obj.minE);
        }
        Ctor.minE = obj.minE;
      }
      if (obj.maxE !== undefined) {
        if (typeof obj.maxE !== 'number' || obj.maxE !== Math.floor(obj.maxE) || obj.maxE < 0 || obj.maxE > 9e15) {
          throw new Error('[decimal.js] DecimalError: Invalid maxE: ' + obj.maxE);
        }
        Ctor.maxE = obj.maxE;
      }
    }
    return {
      precision: Ctor.precision,
      rounding: Ctor.rounding,
      modulo: Ctor.modulo,
      toExpNeg: Ctor.toExpNeg,
      toExpPos: Ctor.toExpPos,
      minE: Ctor.minE,
      maxE: Ctor.maxE
    };
  };

  if (configObj) {
    Ctor.config(configObj);
  }

  // Attach static methods to Ctor
  Ctor.abs = function(x) { return new Ctor(x).abs(); };
  Ctor.acos = function(x) { return new Ctor(x).acos(); };
  Ctor.acosh = function(x) { return new Ctor(x).acosh(); };
  Ctor.add = Ctor.plus = function(x, y) { return new Ctor(x).add(y); };
  Ctor.asin = function(x) { return new Ctor(x).asin(); };
  Ctor.asinh = function(x) { return new Ctor(x).asinh(); };
  Ctor.atan = function(x) { return new Ctor(x).atan(); };
  Ctor.atan2 = function(y, x) { return new Ctor(y).atan2(x); };
  Ctor.atanh = function(x) { return new Ctor(x).atanh(); };
  Ctor.cbrt = function(x) { return new Ctor(x).cbrt(); };
  Ctor.ceil = function(x) { return new Ctor(x).ceil(); };
  Ctor.clamp = Ctor.clampedTo = function(x, min, max) { return new Ctor(x).clamp(min, max); };
  Ctor.clone = Decimal.clone;
  Ctor.cos = function(x) { return new Ctor(x).cos(); };
  Ctor.cosh = function(x) { return new Ctor(x).cosh(); };
  Ctor.div = Ctor.dividedBy = function(x, y) { return new Ctor(x).div(y); };
  Ctor.exp = function(x) { return new Ctor(x).exp(); };
  Ctor.floor = function(x) { return new Ctor(x).floor(); };
  Ctor.hypot = function(...args) {
    let pr = Ctor.precision;
    Ctor.precision = pr + 12;
    let t = new Ctor(0);
    for (let arg of args) {
      let n = new Ctor(arg);
      if (!n.d) {
        if (n.s) {
          Ctor.precision = pr;
          return new Ctor(Infinity);
        }
        t = n;
      } else if (t.d) {
        t = t.plus(n.times(n));
      }
    }
    Ctor.precision = pr;
    return t.sqrt();
  };
  Ctor.ln = function(x) { return new Ctor(x).ln(); };
  Ctor.log = function(x, y) { return new Ctor(x).log(y); };
  Ctor.log10 = function(x) { return new Ctor(x).log(10); };
  Ctor.log2 = function(x) { return new Ctor(x).log(2); };
  Ctor.mod = Ctor.modulo = function(x, y) { return new Ctor(x).mod(y); };
  Ctor.mul = Ctor.times = function(x, y) { return new Ctor(x).mul(y); };
  Ctor.pow = function(x, y) { return new Ctor(x).pow(y); };
  Ctor.round = function(x) { return new Ctor(x).round(); };
  Ctor.sin = function(x) { return new Ctor(x).sin(); };
  Ctor.sinh = function(x) { return new Ctor(x).sinh(); };
  Ctor.sqrt = function(x) { return new Ctor(x).sqrt(); };
  Ctor.sub = Ctor.minus = function(x, y) { return new Ctor(x).minus(y); };
  Ctor.tan = function(x) { return new Ctor(x).tan(); };
  Ctor.tanh = function(x) { return new Ctor(x).tanh(); };
  Ctor.trunc = function(x) { return new Ctor(x).trunc(); };

  return Ctor;
};

// Static methods
Decimal.abs = function(x) { return new Decimal(x).abs(); };
Decimal.acos = function(x) { return new Decimal(x).acos(); };
Decimal.acosh = function(x) { return new Decimal(x).acosh(); };
Decimal.add = Decimal.plus = function(x, y) { return new Decimal(x).add(y); };
Decimal.asin = function(x) { return new Decimal(x).asin(); };
Decimal.asinh = function(x) { return new Decimal(x).asinh(); };
Decimal.atan = function(x) { return new Decimal(x).atan(); };
Decimal.atan2 = function(y, x) { return new Decimal(y).atan2(x); };
Decimal.atanh = function(x) { return new Decimal(x).atanh(); };
Decimal.cbrt = function(x) { return new Decimal(x).cbrt(); };
Decimal.ceil = function(x) { return new Decimal(x).ceil(); };
Decimal.clamp = Decimal.clampedTo = function(x, min, max) { return new Decimal(x).clamp(min, max); };
Decimal.cos = function(x) { return new Decimal(x).cos(); };
Decimal.cosh = function(x) { return new Decimal(x).cosh(); };
Decimal.div = Decimal.dividedBy = function(x, y) { return new Decimal(x).div(y); };
Decimal.exp = function(x) { return new Decimal(x).exp(); };
Decimal.floor = function(x) { return new Decimal(x).floor(); };
Decimal.hypot = function (...args) {
  const strArgs = args.map(arg => (arg instanceof Decimal ? arg.valueOf() : arg));
  return wrapRes(callGo('hypot', strArgs, this));
};
Decimal.ln = function(x) { return new Decimal(x).ln(); };
Decimal.log = function(x, y) { return new Decimal(x).log(y); };
Decimal.log10 = function(x) { return new Decimal(x).log(10); };
Decimal.log2 = function(x) { return new Decimal(x).log(2); };
Decimal.max = function(...args) {
  if (args.length === 0) return new Decimal(NaN);
  let max = new Decimal(args[0]);
  for (let i = 1; i < args.length; i++) {
    let current = new Decimal(args[i]);
    if (max.isNaN() || current.isNaN()) return new Decimal(NaN);
    if (current.gt(max) || (current.eq(max) && current.s === 1 && max.s === -1)) max = current;
  }
  return max;
};
Decimal.min = function(...args) {
  if (args.length === 0) return new Decimal(NaN);
  let min = new Decimal(args[0]);
  for (let i = 1; i < args.length; i++) {
    let current = new Decimal(args[i]);
    if (min.isNaN() || current.isNaN()) return new Decimal(NaN);
    if (current.lt(min) || (current.eq(min) && current.s === -1 && min.s === 1)) min = current;
  }
  return min;
};
Decimal.mod = Decimal.modulo = function(x, y) { return new Decimal(x).mod(y); };
Decimal.mul = Decimal.times = function(x, y) { return new Decimal(x).times(y); };
Decimal.pow = function(x, y) { return new Decimal(x).pow(y); };
Decimal.random = function (dp) {
  const Ctor = this || Decimal;
  if (dp !== undefined) {
    checkInt32(dp, 0, MAX_DIGITS);
    const oldPrec = Ctor.precision;
    Ctor.precision = dp;
    const r = wrapJS(DecimalJS.random(dp), Ctor);
    Ctor.precision = oldPrec;
    return r;
  }
  return wrapJS(DecimalJS.random(), Ctor);
};
Decimal.round = function(x) { return new Decimal(x).round(); };
Decimal.sign = function (x) {
  let d = new Decimal(x);
  if (d.isNaN()) return NaN;
  if (d.isZero()) return d.s < 0 ? -0 : 0;
  return d.s;
};
Decimal.sin = function(x) { return new Decimal(x).sin(); };
Decimal.sinh = function(x) { return new Decimal(x).sinh(); };
Decimal.sqrt = function(x) { return new Decimal(x).sqrt(); };
Decimal.sub = Decimal.minus = function(x, y) { return new Decimal(x).minus(y); };
Decimal.sum = function (...args) {
  if (args.length === 0) return new Decimal(NaN);
  let total = new Decimal(0);
  for (let arg of args) {
    total = total.plus(arg);
  }
  return total;
};
Decimal.tan = function(x) { return new Decimal(x).tan(); };
Decimal.tanh = function(x) { return new Decimal(x).tanh(); };
Decimal.trunc = function(x) { return new Decimal(x).trunc(); };
Decimal.isDecimal = Decimal.isDecimalInstance = function (obj) {
  return obj instanceof Decimal || (obj && typeof obj === 'object' && obj.d !== undefined && obj.s !== undefined && obj.e !== undefined);
};

// Instance methods
Decimal.prototype.plus = Decimal.prototype.add = function (y) {
  y = new this.constructor(y);
  return wrapRes(callGo('plus', [this.valueOf(), y.valueOf()], this.constructor));
};

Decimal.prototype.minus = Decimal.prototype.sub = function (y) {
  y = new this.constructor(y);
  return wrapRes(callGo('minus', [this.valueOf(), y.valueOf()], this.constructor));
};

Decimal.prototype.times = Decimal.prototype.mul = function (y) {
  y = new this.constructor(y);
  return wrapRes(callGo('times', [this.valueOf(), y.valueOf()], this.constructor));
};

Decimal.prototype.dividedBy = Decimal.prototype.div = function (y) {
  y = new this.constructor(y);
  return wrapRes(callGo('dividedBy', [this.valueOf(), y.valueOf()], this.constructor));
};

Decimal.prototype.modulo = Decimal.prototype.mod = function (y) {
  y = new this.constructor(y);
  return wrapRes(callGo('mod', [this.valueOf(), y.valueOf()], this.constructor));
};

Decimal.prototype.squareRoot = Decimal.prototype.sqrt = function () {
  return wrapRes(callGo('sqrt', [this.valueOf()], this.constructor));
};

Decimal.prototype.cubeRoot = Decimal.prototype.cbrt = function () {
  return wrapRes(callGo('cbrt', [this.valueOf()], this.constructor));
};

Decimal.prototype.toPower = Decimal.prototype.pow = function (y) {
  y = new this.constructor(y);
  return wrapRes(callGo('pow', [this.valueOf(), y.valueOf()], this.constructor));
};

Decimal.prototype.naturalLogarithm = Decimal.prototype.ln = function () {
  return wrapRes(callGo('ln', [this.valueOf()], this.constructor));
};

Decimal.prototype.naturalExponential = Decimal.prototype.exp = function () {
  return wrapRes(callGo('exp', [this.valueOf()], this.constructor));
};

Decimal.prototype.sine = Decimal.prototype.sin = function () {
  return wrapRes(callGo('sin', [this.valueOf()], this.constructor));
};

Decimal.prototype.cosine = Decimal.prototype.cos = function () {
  return wrapRes(callGo('cos', [this.valueOf()], this.constructor));
};

Decimal.prototype.tangent = Decimal.prototype.tan = function () {
  return wrapRes(callGo('tan', [this.valueOf()], this.constructor));
};

Decimal.prototype.arcsine = Decimal.prototype.inverseSine = Decimal.prototype.asin = function () {
  return wrapRes(callGo('asin', [this.valueOf()], this.constructor));
};

Decimal.prototype.arccosine = Decimal.prototype.inverseCosine = Decimal.prototype.acos = function () {
  return wrapRes(callGo('acos', [this.valueOf()], this.constructor));
};

Decimal.prototype.arctangent = Decimal.prototype.inverseTangent = Decimal.prototype.atan = function () {
  return wrapRes(callGo('atan', [this.valueOf()], this.constructor));
};

Decimal.prototype.hyperbolicSine = Decimal.prototype.sinh = function () {
  return wrapRes(callGo('sinh', [this.valueOf()], this.constructor));
};

Decimal.prototype.hyperbolicCosine = Decimal.prototype.cosh = function () {
  return wrapRes(callGo('cosh', [this.valueOf()], this.constructor));
};

Decimal.prototype.hyperbolicTangent = Decimal.prototype.tanh = function () {
  return wrapRes(callGo('tanh', [this.valueOf()], this.constructor));
};

Decimal.prototype.inverseHyperbolicSine = Decimal.prototype.asinh = function () {
  return wrapRes(callGo('asinh', [this.valueOf()], this.constructor));
};

Decimal.prototype.inverseHyperbolicCosine = Decimal.prototype.acosh = function () {
  return wrapRes(callGo('acosh', [this.valueOf()], this.constructor));
};

Decimal.prototype.inverseHyperbolicTangent = Decimal.prototype.atanh = function () {
  return wrapRes(callGo('atanh', [this.valueOf()], this.constructor));
};

Decimal.prototype.arctan2 = Decimal.prototype.atan2 = function (y) {
  return wrapRes(callGo('atan2', [this.valueOf(), y instanceof Decimal ? y.valueOf() : y], this.constructor));
};

Decimal.prototype.abs = Decimal.prototype.absoluteValue = function () {
  return wrapRes(callGo('abs', [this.valueOf()], this.constructor));
};

Decimal.prototype.neg = Decimal.prototype.negated = function () {
  return wrapRes(callGo('neg', [this.valueOf()], this.constructor));
};

Decimal.prototype.trunc = Decimal.prototype.truncated = function () {
  return wrapRes(callGo('trunc', [this.valueOf()], this.constructor));
};

Decimal.prototype.floor = function () {
  return wrapRes(callGo('floor', [this.valueOf()], this.constructor));
};

Decimal.prototype.ceil = function () {
  return wrapRes(callGo('ceil', [this.valueOf()], this.constructor));
};

Decimal.prototype.round = function () {
  return wrapRes(callGo('round', [this.valueOf()], this.constructor));
};

Decimal.prototype.cmp = Decimal.prototype.comparedTo = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('cmp', [this.valueOf(), yVal], this.constructor);
  return res.isCmp ? res.cmpRes : NaN;
};

Decimal.prototype.eq = Decimal.prototype.equals = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('eq', [this.valueOf(), yVal], this.constructor);
  return res.isBool ? res.boolRes : false;
};

Decimal.prototype.gt = Decimal.prototype.greaterThan = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('gt', [this.valueOf(), yVal], this.constructor);
  return res.isBool ? res.boolRes : false;
};

Decimal.prototype.gte = Decimal.prototype.greaterThanOrEqualTo = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('gte', [this.valueOf(), yVal], this.constructor);
  return res.isBool ? res.boolRes : false;
};

Decimal.prototype.lt = Decimal.prototype.lessThan = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('lt', [this.valueOf(), yVal], this.constructor);
  return res.isBool ? res.boolRes : false;
};

Decimal.prototype.lte = Decimal.prototype.lessThanOrEqualTo = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const res = callGo('lte', [this.valueOf(), yVal], this.constructor);
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

Decimal.prototype.toString = function () {
  if (this.d && this.d.length === 1 && this.d[0] === 0 && this.s === -1) {
    return '0';
  }
  return this._str || '0';
};

Decimal.prototype.toJSON = Decimal.prototype.valueOf = function () {
  if (this.d && this.d.length === 1 && this.d[0] === 0 && this.s === -1) {
    return '-0';
  }
  return this._str || '0';
};

// New delegated prototype methods
// New delegated prototype methods
Decimal.prototype.clamp = Decimal.prototype.clampedTo = function (min, max) {
  const minVal = min instanceof Decimal ? min.valueOf() : min;
  const maxVal = max instanceof Decimal ? max.valueOf() : max;
  return wrapRes(callGo('clamp', [this.valueOf(), minVal, maxVal], this.constructor));
};

Decimal.prototype.divToInt = Decimal.prototype.dividedToIntegerBy = function (y) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  return wrapRes(callGo('divToInt', [this.valueOf(), yVal], this.constructor));
};

Decimal.prototype.dp = Decimal.prototype.decimalPlaces = function () {
  const res = callGo('dp', [this.valueOf()], this.constructor);
  const val = parseInt(res.string || res.str, 10);
  return val < 0 ? NaN : val;
};

Decimal.prototype.sd = Decimal.prototype.precision = function (z) {
  if (z !== undefined && z !== true && z !== false && z !== 1 && z !== 0) {
    throw new Error('[decimal.js] DecimalError: Invalid argument: ' + z);
  }
  const res = callGo('sd', [this.valueOf(), z], this.constructor);
  const val = parseInt(res.string || res.str, 10);
  return val < 0 ? NaN : val;
};

Decimal.prototype.isNeg = Decimal.prototype.isNegative = function () {
  return this.s < 0;
};

Decimal.prototype.isPos = Decimal.prototype.isPositive = function () {
  return this.s > 0;
};

Decimal.prototype.isInt = Decimal.prototype.isInteger = function () {
  const res = callGo('isInt', [this.valueOf()], this.constructor);
  return (res.string || res.str) === 'true';
};

Decimal.prototype.toFixed = function (dp, rm) {
  if (dp !== undefined) checkInt32(dp, 0, MAX_DIGITS);
  if (rm !== undefined) checkInt32(rm, 0, 8);
  const oldRounding = this.constructor.rounding;
  if (rm !== undefined) this.constructor.rounding = rm;
  const res = callGo('toFixed', [this.valueOf(), dp], this.constructor);
  if (rm !== undefined) this.constructor.rounding = oldRounding;
  return res.string || res.str;
};

Decimal.prototype.toExponential = function (dp, rm) {
  if (dp !== undefined) checkInt32(dp, 0, MAX_DIGITS);
  if (rm !== undefined) checkInt32(rm, 0, 8);
  if (dp === undefined) dp = -1;
  const oldRounding = this.constructor.rounding;
  if (rm !== undefined) this.constructor.rounding = rm;
  const res = callGo('toExponential', [this.valueOf(), dp], this.constructor);
  if (rm !== undefined) this.constructor.rounding = oldRounding;
  return res.string || res.str;
};

Decimal.prototype.toPrecision = function (sd, rm) {
  if (sd !== undefined) checkInt32(sd, 1, MAX_DIGITS);
  if (rm !== undefined) checkInt32(rm, 0, 8);
  if (sd === undefined) return this.toString();
  const oldRounding = this.constructor.rounding;
  if (rm !== undefined) this.constructor.rounding = rm;
  const res = callGo('toPrecision', [this.valueOf(), sd], this.constructor);
  if (rm !== undefined) this.constructor.rounding = oldRounding;
  return res.string || res.str;
};

Decimal.prototype.toBinary = function (sd, rm) {
  const res = callGo('toBinary', [this.valueOf(), sd, rm], this.constructor);
  return res.string || res.str;
};

Decimal.prototype.toHex = Decimal.prototype.toHexadecimal = function (sd, rm) {
  const res = callGo('toHex', [this.valueOf(), sd, rm], this.constructor);
  return res.string || res.str;
};

Decimal.prototype.toOctal = function (sd, rm) {
  const res = callGo('toOctal', [this.valueOf(), sd, rm], this.constructor);
  return res.string || res.str;
};

Decimal.prototype.toDP = Decimal.prototype.toDecimalPlaces = function (dp, rm) {
  if (dp === undefined) {
    return this;
  }
  checkInt32(dp, 0, MAX_DIGITS);
  if (rm !== undefined) checkInt32(rm, 0, 8);
  const rmVal = rm !== undefined ? rm : this.constructor.rounding;
  return wrapRes(callGo('toDP', [this.valueOf(), dp, rmVal], this.constructor));
};

Decimal.prototype.toSD = Decimal.prototype.toSignificantDigits = function (sd, rm) {
  if (sd !== undefined) {
    checkInt32(sd, 1, MAX_DIGITS);
  }
  if (rm !== undefined) {
    checkInt32(rm, 0, 8);
  }
  const sdVal = sd !== undefined ? sd : this.constructor.precision;
  const rmVal = rm !== undefined ? rm : this.constructor.rounding;
  return wrapRes(callGo('toSD', [this.valueOf(), sdVal, rmVal], this.constructor));
};

Decimal.prototype.toFraction = function (maxD) {
  const maxDVal = maxD !== undefined ? (maxD instanceof Decimal ? maxD.valueOf() : maxD) : undefined;
  const res = callGo('toFraction', [this.valueOf(), maxDVal], this.constructor);
  const str = res.string || res.str || '0/1';
  const parts = str.split('/');
  const Ctor = this.constructor;
  return [new Ctor(parts[0]), new Ctor(parts[1] || '1')];
};

Decimal.prototype.toNearest = function (y, rm) {
  const yVal = y instanceof Decimal ? y.valueOf() : y;
  const rmVal = rm !== undefined ? rm : this.constructor.rounding;
  return wrapRes(callGo('toNearest', [this.valueOf(), yVal, rmVal], this.constructor));
};

Decimal.prototype.toNumber = function () {
  return Number(this.valueOf());
};

Decimal.prototype.logarithm = Decimal.prototype.log = function (base) {
  const baseVal = base !== undefined ? (base instanceof Decimal ? base.valueOf() : base) : undefined;
  return wrapRes(callGo('log', [this.valueOf(), baseVal], this.constructor));
};

module.exports = Decimal;
