package decimal

// Sqrt computes square root sqrt(x) using context c settings.
// Matches decimal.js squareRoot / sqrt (lines 1432-1522).
func (c *Context) Sqrt(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() || x.s < 0 {
		return &Decimal{s: 0}
	}
	if x.IsZero() || x.d == nil {
		return &Decimal{s: x.s, e: x.e, d: x.d}
	}
	one, _ := c.New(1)
	if x.Eq(one) {
		return one
	}

	wpr := c.Precision + 10
	// Newton-Raphson initial seed: r_0 = x / 2
	half, _ := c.New(0.5)
	r := c.Mul(x, half)

	for i := 0; i < 50; i++ {
		// r_{k+1} = 0.5 * (r_k + x / r_k)
		xDivR := c.divide(x, r, wpr, RoundDown, false)
		sum := c.Add(r, xDivR)
		nextR := c.Mul(half, sum)
		if nextR.Eq(r) {
			break
		}
		r = nextR
	}

	return c.finalise(r, c.Precision, c.Rounding, false)
}

// Sqrt computes square root sqrt(x) using default context.
func (x *Decimal) Sqrt() *Decimal {
	return globalContext.Sqrt(x)
}

// SquareRoot is an alias for Sqrt.
func (x *Decimal) SquareRoot() *Decimal {
	return x.Sqrt()
}

// Cbrt computes cube root cbrt(x) using context c settings.
// Matches decimal.js cubeRoot / cbrt (lines 334-421).
func (c *Context) Cbrt(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if x.IsZero() || x.d == nil {
		return &Decimal{s: x.s, e: x.e, d: x.d}
	}

	wpr := c.Precision + 10
	r, _ := c.New(x)
	two, _ := c.New(2)

	// Halley's cubic method: r_{k+1} = r_k * (r_k^3 + 2X) / (2r_k^3 + X)
	for i := 0; i < 50; i++ {
		r2 := c.Mul(r, r)
		r3 := c.Mul(r2, r)
		twoX := c.Mul(two, x)
		num := c.Add(r3, twoX)

		twoR3 := c.Mul(two, r3)
		den := c.Add(twoR3, x)

		frac := c.divide(num, den, wpr, RoundDown, false)
		nextR := c.Mul(r, frac)
		if nextR.Eq(r) {
			break
		}
		r = nextR
	}

	return c.finalise(r, c.Precision, c.Rounding, false)
}

// Cbrt computes cube root cbrt(x) using default context.
func (x *Decimal) Cbrt() *Decimal {
	return globalContext.Cbrt(x)
}

// CubeRoot is an alias for Cbrt.
func (x *Decimal) CubeRoot() *Decimal {
	return x.Cbrt()
}
