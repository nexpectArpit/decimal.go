package decimal

// Mod computes x % y using context c settings.
// Matches decimal.js modulo / mod (lines 1133-1185).
func (c *Context) Mod(x, y *Decimal) *Decimal {
	if x == nil || y == nil || !x.IsFinite() || !y.IsFinite() || y.IsZero() {
		return &Decimal{s: 0}
	}
	if x.IsZero() {
		return &Decimal{s: x.s, e: 0, d: []int32{0}}
	}

	// Compute quotient q = trunc(x / y)
	q := c.divide(x, y, 0, RoundDown, true)
	qTrunc := c.finalise(q, int(q.e)+1, RoundDown, false)

	// r = x - (y * qTrunc)
	prod := c.Mul(y, qTrunc)
	rem := c.Sub(x, prod)

	return rem
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
