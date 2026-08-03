package decimal

// Mod computes x % y using context c settings.
// Matches decimal.js modulo / mod (lines 1438-1471).
func (c *Context) Mod(x, y *Decimal) *Decimal {
	if x == nil || y == nil {
		return &Decimal{s: 0}
	}

	// Return NaN if x is ±Infinity or NaN, or y is NaN or ±0.
	if x.d == nil || y.s == 0 || (len(y.d) > 0 && y.d[0] == 0) {
		return &Decimal{s: 0}
	}

	// Return x if y is ±Infinity or x is ±0.
	if y.d == nil || (len(x.d) > 0 && x.d[0] == 0) {
		return c.finalise(new(Decimal).Set(x), c.Precision, c.Rounding, false)
	}

	// Compute required precision for intermediate operations to avoid truncation of remainder.
	wpr := c.Precision + 40
	diffE := int(x.e - y.e)
	if diffE > 0 {
		wpr += diffE
	}

	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	var q *Decimal
	if c.Modulo == ModuloEuclid {
		// Euclidian division: q = sign(y) * floor(x / abs(y))
		// result = x - q * y    where  0 <= result < abs(y)
		q = evalCtx.divide(x, y.Abs(), 0, RoundingMode(ModuloFloor), true)
		q.s *= y.s
	} else {
		q = evalCtx.divide(x, y, 0, RoundingMode(c.Modulo), true)
	}

	prod := evalCtx.Mul(q, y)
	rem := evalCtx.Sub(x, prod)

	return c.finalise(rem, c.Precision, c.Rounding, false)
}

// Mod computes x % y using default context.
func (x *Decimal) Mod(y *Decimal) *Decimal {
	return globalContext.Mod(x, y)
}

// Modulo is an alias for Mod.
func (x *Decimal) Modulo(y *Decimal) *Decimal {
	return x.Mod(y)
}

// DivToInt computes integer quotient of x / y using context c settings.
// Matches decimal.js dividedToIntegerBy / divToInt (lines 478-482).
func (c *Context) DivToInt(x, y *Decimal) *Decimal {
	q := c.divide(x, y, 0, RoundDown, true)
	return c.finalise(q, c.Precision, c.Rounding, false)
}

// DivToInt computes integer quotient of x / y using default context.
func (x *Decimal) DivToInt(y *Decimal) *Decimal {
	return globalContext.DivToInt(x, y)
}

// DividedToIntegerBy is an alias for DivToInt.
func (x *Decimal) DividedToIntegerBy(y *Decimal) *Decimal {
	return x.DivToInt(y)
}
