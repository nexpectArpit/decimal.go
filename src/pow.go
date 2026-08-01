package decimal

// intPow performs binary exponentiation x^n for small integer exponents n.
// Matches decimal.js intPow() (lines 3208-3241).
func (c *Context) intPow(x *Decimal, n int64) *Decimal {
	res, _ := c.New(1)
	base, _ := c.New(x)

	if n < 0 {
		n = -n
		base = c.Div(&Decimal{s: 1, e: 0, d: []int32{1}}, base)
	}

	for n > 0 {
		if n%2 == 1 {
			res = c.Mul(res, base)
		}
		base = c.Mul(base, base)
		n /= 2
	}

	return res
}

// Pow computes x^y using context c settings.
// Matches decimal.js toPower / pow (lines 2268-2365).
func (c *Context) Pow(x, y *Decimal) *Decimal {
	if x == nil || y == nil || x.IsNaN() || y.IsNaN() {
		return &Decimal{s: 0}
	}

	// pow(x, ±0) = 1
	if y.IsZero() {
		return &Decimal{s: 1, e: 0, d: []int32{1}}
	}

	// Check if y is small integer for fast binary exponentiation
	if y.IsInt() && mathAbs(y.e) < 15 {
		n := int64(0)
		if len(y.d) > 0 {
			n = int64(y.d[0])
			if y.s < 0 {
				n = -n
			}
		}
		if n >= -MaxSafeInteger && n <= MaxSafeInteger {
			return c.intPow(x, n)
		}
	}

	// x^y = exp(y * ln(x))
	lnX := c.Ln(x)
	yLnX := c.Mul(y, lnX)
	return c.Exp(yLnX)
}

// Pow computes x^y using default context.
func (x *Decimal) Pow(y *Decimal) *Decimal {
	return globalContext.Pow(x, y)
}

// ToPower is an alias for Pow.
func (x *Decimal) ToPower(y *Decimal) *Decimal {
	return x.Pow(y)
}

func mathAbs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
