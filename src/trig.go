package decimal

// Sin computes sine in radians using context c settings.
// Matches decimal.js sine / sin (lines 1398-1419, 3687-3719).
func (c *Context) Sin(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() || x.d == nil {
		return &Decimal{s: 0} // sin(Inf/NaN) = NaN
	}
	if x.IsZero() {
		return &Decimal{s: x.s, e: 0, d: []int32{0}} // sin(0) = 0
	}

	wpr := c.Precision + 10
	sum, _ := c.New(x)
	pow, _ := c.New(x)
	den, _ := c.New(1)
	x2 := c.finalise(c.Mul(x, x), wpr, RoundDown, false)
	sign := int8(1)

	for i := 1; i < 150; i += 2 {
		if i > 1 {
			pow = c.finalise(c.Mul(pow, x2), wpr, RoundDown, false)
			i1, _ := c.New(int64(i - 1))
			i2, _ := c.New(int64(i))
			den = c.Mul(den, c.Mul(i1, i2))
			term := c.divide(pow, den, wpr, RoundDown, false)
			sign = -sign
			if sign < 0 {
				sum = c.Sub(sum, term)
			} else {
				sum = c.Add(sum, term)
			}
			if len(term.d) == 0 || term.d[0] == 0 {
				break
			}
		}
	}

	return c.finalise(sum, c.Precision, c.Rounding, false)
}

// Sin computes sine in radians using default context.
func (x *Decimal) Sin() *Decimal {
	return globalContext.Sin(x)
}

// Sine is an alias for Sin.
func (x *Decimal) Sine() *Decimal {
	return x.Sin()
}

// Cos computes cosine in radians using context c settings.
// Matches decimal.js cosine / cos (lines 294-315, 2641-2672).
func (c *Context) Cos(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() || x.d == nil {
		return &Decimal{s: 0}
	}
	if x.IsZero() {
		return &Decimal{s: 1, e: 0, d: []int32{1}} // cos(0) = 1
	}

	wpr := c.Precision + 10
	sum, _ := c.New(1)
	pow, _ := c.New(1)
	den, _ := c.New(1)
	x2 := c.finalise(c.Mul(x, x), wpr, RoundDown, false)
	sign := int8(1)

	for i := 2; i < 150; i += 2 {
		pow = c.finalise(c.Mul(pow, x2), wpr, RoundDown, false)
		i1, _ := c.New(int64(i - 1))
		i2, _ := c.New(int64(i))
		den = c.Mul(den, c.Mul(i1, i2))
		term := c.divide(pow, den, wpr, RoundDown, false)
		sign = -sign
		if sign < 0 {
			sum = c.Sub(sum, term)
		} else {
			sum = c.Add(sum, term)
		}
		if len(term.d) == 0 || term.d[0] == 0 {
			break
		}
	}

	return c.finalise(sum, c.Precision, c.Rounding, false)
}

// Cos computes cosine in radians using default context.
func (x *Decimal) Cos() *Decimal {
	return globalContext.Cos(x)
}

// Cosine is an alias for Cos.
func (x *Decimal) Cosine() *Decimal {
	return x.Cos()
}

// Tan computes tangent in radians (sin(x)/cos(x)) using context c settings.
func (c *Context) Tan(x *Decimal) *Decimal {
	sinX := c.Sin(x)
	cosX := c.Cos(x)
	return c.Div(sinX, cosX)
}

// Tan computes tangent in radians using default context.
func (x *Decimal) Tan() *Decimal {
	return globalContext.Tan(x)
}

// Tangent is an alias for Tan.
func (x *Decimal) Tangent() *Decimal {
	return x.Tan()
}

// Asin computes inverse sine (asin) using context c settings.
func (c *Context) Asin(x *Decimal) *Decimal {
	one, _ := c.New(1)
	// asin(x) = atan(x / sqrt(1 - x^2))
	x2 := c.Mul(x, x)
	oneMinusX2 := c.Sub(one, x2)
	sqrtVal := c.Sqrt(oneMinusX2)
	frac := c.Div(x, sqrtVal)
	return c.Atan(frac)
}

// Asin computes inverse sine using default context.
func (x *Decimal) Asin() *Decimal {
	return globalContext.Asin(x)
}

// Acos computes inverse cosine (acos) using context c settings.
func (c *Context) Acos(x *Decimal) *Decimal {
	// acos(x) = pi/2 - asin(x)
	halfPi := c.getPi(c.Precision + 5)
	two, _ := c.New(2)
	halfPi = c.Div(halfPi, two)
	asinX := c.Asin(x)
	return c.Sub(halfPi, asinX)
}

// Acos computes inverse cosine using default context.
func (x *Decimal) Acos() *Decimal {
	return globalContext.Acos(x)
}

// Atan computes inverse tangent (atan) using context c settings.
func (c *Context) Atan(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if x.IsZero() {
		return &Decimal{s: x.s, e: 0, d: []int32{0}}
	}

	wpr := c.Precision + 10
	sum, _ := c.New(x)
	pow, _ := c.New(x)
	x2 := c.finalise(c.Mul(x, x), wpr, RoundDown, false)
	sign := int8(1)

	for i := 3; i < 200; i += 2 {
		pow = c.finalise(c.Mul(pow, x2), wpr, RoundDown, false)
		iDec, _ := c.New(int64(i))
		term := c.divide(pow, iDec, wpr, RoundDown, false)
		sign = -sign
		if sign < 0 {
			sum = c.Sub(sum, term)
		} else {
			sum = c.Add(sum, term)
		}
		if len(term.d) == 0 || term.d[0] == 0 {
			break
		}
	}

	return c.finalise(sum, c.Precision, c.Rounding, false)
}

// Atan computes inverse tangent using default context.
func (x *Decimal) Atan() *Decimal {
	return globalContext.Atan(x)
}
