package decimal

// Sinh computes hyperbolic sine sinh(x) = (exp(x) - exp(-x)) / 2 using context c settings.
func (c *Context) Sinh(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if !x.IsFinite() || x.IsZero() {
		return &Decimal{s: x.s, e: x.e, d: append([]int32(nil), x.d...)}
	}

	sd := len(x.d) * 7
	absE := x.e
	if absE < 0 {
		absE = -absE
	}
	maxESd := absE
	if int64(sd) > maxESd {
		maxESd = int64(sd)
	}
	wpr := int(int64(c.Precision) + maxESd + 40)

	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	expX := evalCtx.Exp(x)
	negX := x.Neg()
	expNegX := evalCtx.Exp(negX)
	diff := evalCtx.Sub(expX, expNegX)
	two, _ := evalCtx.New(2)
	res := evalCtx.Div(diff, two)

	return c.finalise(res, c.Precision, c.Rounding, false)
}

// Sinh computes hyperbolic sine using default context.
func (x *Decimal) Sinh() *Decimal {
	return globalContext.Sinh(x)
}

// HyperbolicSine is an alias for Sinh.
func (x *Decimal) HyperbolicSine() *Decimal {
	return x.Sinh()
}

// Cosh computes hyperbolic cosine cosh(x) = (exp(x) + exp(-x)) / 2 using context c settings.
func (c *Context) Cosh(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if !x.IsFinite() {
		return &Decimal{s: 1, e: 0, d: nil} // Infinity
	}
	if x.IsZero() {
		return &Decimal{s: 1, e: 0, d: []int32{1}} // 1
	}

	sd := len(x.d) * 7
	absE := x.e
	if absE < 0 {
		absE = -absE
	}
	maxESd := absE
	if int64(sd) > maxESd {
		maxESd = int64(sd)
	}
	wpr := int(int64(c.Precision) + maxESd + 40)

	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	expX := evalCtx.Exp(x)
	negX := x.Neg()
	expNegX := evalCtx.Exp(negX)
	sum := evalCtx.Add(expX, expNegX)
	two, _ := evalCtx.New(2)
	res := evalCtx.Div(sum, two)

	return c.finalise(res, c.Precision, c.Rounding, false)
}

// Cosh computes hyperbolic cosine using default context.
func (x *Decimal) Cosh() *Decimal {
	return globalContext.Cosh(x)
}

// HyperbolicCosine is an alias for Cosh.
func (x *Decimal) HyperbolicCosine() *Decimal {
	return x.Cosh()
}

// Tanh computes hyperbolic tangent tanh(x) = sinh(x) / cosh(x) using context c settings.
func (c *Context) Tanh(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if !x.IsFinite() {
		return &Decimal{s: x.s, e: 0, d: []int32{1}} // ±1
	}
	if x.IsZero() {
		return &Decimal{s: x.s, e: 0, d: []int32{0}}
	}

	wpr := c.Precision + 40

	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	sinhX := evalCtx.Sinh(x)
	coshX := evalCtx.Cosh(x)
	res := evalCtx.Div(sinhX, coshX)

	return c.finalise(res, c.Precision, c.Rounding, false)
}

// Tanh computes hyperbolic tangent using default context.
func (x *Decimal) Tanh() *Decimal {
	return globalContext.Tanh(x)
}

// HyperbolicTangent is an alias for Tanh.
func (x *Decimal) HyperbolicTangent() *Decimal {
	return x.Tanh()
}

// Asinh computes inverse hyperbolic sine asinh(x) = ln(x + sqrt(x^2 + 1)) using context c settings.
func (c *Context) Asinh(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if !x.IsFinite() || x.IsZero() {
		return &Decimal{s: x.s, e: x.e, d: append([]int32(nil), x.d...)}
	}

	sd := len(x.d) * 7
	absE := x.e
	if absE < 0 {
		absE = -absE
	}
	maxESd := absE
	if int64(sd) > maxESd {
		maxESd = int64(sd)
	}
	wpr := int(int64(c.Precision) + 2*maxESd + 40)

	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	one, _ := evalCtx.New(1)
	x2 := evalCtx.Mul(x, x)
	x2Plus1 := evalCtx.Add(x2, one)
	sqrtVal := evalCtx.Sqrt(x2Plus1)
	sum := evalCtx.Add(x, sqrtVal)
	res := evalCtx.Ln(sum)

	return c.finalise(res, c.Precision, c.Rounding, false)
}

// Asinh computes inverse hyperbolic sine using default context.
func (x *Decimal) Asinh() *Decimal {
	return globalContext.Asinh(x)
}

// Acosh computes inverse hyperbolic cosine acosh(x) = ln(x + sqrt(x^2 - 1)) using context c settings.
func (c *Context) Acosh(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	one, _ := c.New(1)
	if x.Lt(one) {
		return &Decimal{s: 0} // NaN for x < 1
	}
	if x.Eq(one) {
		return &Decimal{s: 1, e: 0, d: []int32{0}} // 0 for x == 1
	}
	if !x.IsFinite() {
		return &Decimal{s: x.s, e: 0, d: nil} // Infinity
	}

	sd := len(x.d) * 7
	absE := x.e
	if absE < 0 {
		absE = -absE
	}
	maxESd := absE
	if int64(sd) > maxESd {
		maxESd = int64(sd)
	}
	wpr := int(int64(c.Precision) + absE + maxESd + 40)

	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown
	x2 := evalCtx.Mul(x, x)
	x2Minus1 := evalCtx.Sub(x2, one)
	sqrtVal := evalCtx.Sqrt(x2Minus1)
	sum := evalCtx.Add(x, sqrtVal)
	res := evalCtx.Ln(sum)

	return c.finalise(res, c.Precision, c.Rounding, false)
}

// Acosh computes inverse hyperbolic cosine using default context.
func (x *Decimal) Acosh() *Decimal {
	return globalContext.Acosh(x)
}

// Atanh computes inverse hyperbolic tangent atanh(x) = 0.5 * ln((1 + x) / (1 - x)) using context c settings.
func (c *Context) Atanh(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if !x.IsFinite() {
		return &Decimal{s: 0} // NaN
	}
	if x.IsZero() {
		return &Decimal{s: x.s, e: 0, d: []int32{0}}
	}
	
	one, _ := c.New(1)
	absX := x.Abs()
	if absX.Eq(one) {
		return &Decimal{s: x.s, e: 0, d: nil} // ±Infinity
	}
	if absX.Gt(one) {
		return &Decimal{s: 0} // NaN
	}

	sd := len(x.d) * 7
	xsd := sd
	pr := c.Precision

	maxVal := xsd
	if pr > maxVal {
		maxVal = pr
	}
	if int64(maxVal) < 2*(-x.e)-1 {
		return c.finalise(new(Decimal).Set(x), pr, c.Rounding, false)
	}

	wpr := xsd - int(x.e)
	
	evalCtx := c.Clone()
	evalCtx.Precision = wpr + pr + 40
	evalCtx.Rounding = RoundDown
	
	num := evalCtx.Add(x, one)
	den := evalCtx.Sub(one, x)
	frac := evalCtx.divide(num, den, wpr+pr+40, RoundDown, false)

	evalCtx2 := c.Clone()
	evalCtx2.Precision = pr + 40
	evalCtx2.Rounding = RoundDown
	
	lnFrac := evalCtx2.Ln(frac)
	half, _ := evalCtx2.New(0.5)
	res := evalCtx2.Mul(half, lnFrac)

	return c.finalise(res, c.Precision, c.Rounding, false)
}

// Atanh computes inverse hyperbolic tangent using default context.
func (x *Decimal) Atanh() *Decimal {
	return globalContext.Atanh(x)
}
