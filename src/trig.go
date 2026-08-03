package decimal

func (c *Context) toLessThanHalfPi(x *Decimal) (*Decimal, int) {
	wpr := c.Precision + 12
	absE := x.e
	if absE < 0 {
		absE = -absE
	}
	if int(absE) > 0 {
		wpr += int(absE)
	}
	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	isNeg := x.s < 0
	pi := evalCtx.getPi(wpr)
	half, _ := evalCtx.New(0.5)
	halfPi := evalCtx.Mul(pi, half)

	xAbs := x.Abs()

	if xAbs.Lte(halfPi) {
		quadrant := 1
		if isNeg {
			quadrant = 4
		}
		return xAbs, quadrant
	}

	t := evalCtx.DivToInt(xAbs, pi)

	var quadrant int
	if t.IsZero() {
		quadrant = 2
		if isNeg {
			quadrant = 3
		}
	} else {
		xAbs = evalCtx.Sub(xAbs, evalCtx.Mul(t, pi))

		// 0 <= xAbs < pi
		if xAbs.Lte(halfPi) {
			two, _ := evalCtx.New(2)
			isOddT := !evalCtx.Mod(t, two).IsZero()

			if isOddT {
				if isNeg {
					quadrant = 2
				} else {
					quadrant = 3
				}
			} else {
				if isNeg {
					quadrant = 4
				} else {
					quadrant = 1
				}
			}
			return xAbs, quadrant
		}

		two, _ := evalCtx.New(2)
		isOddT := !evalCtx.Mod(t, two).IsZero()

		if isOddT {
			if isNeg {
				quadrant = 1
			} else {
				quadrant = 4
			}
		} else {
			if isNeg {
				quadrant = 3
			} else {
				quadrant = 2
			}
		}
	}

	res := evalCtx.Sub(xAbs, pi).Abs()
	return res, quadrant
}

// Sin computes sine in radians using context c settings.
func (c *Context) Sin(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() || x.d == nil {
		return &Decimal{s: 0} // sin(Inf/NaN) = NaN
	}
	if x.IsZero() {
		return &Decimal{s: x.s, e: 0, d: []int32{0}} // sin(0) = 0
	}

	wpr := c.Precision + 20
	absE := x.e
	if absE < 0 {
		absE = -absE
	}
	if int(absE) > 0 {
		wpr += int(absE)
	}
	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	xReduced, quadrant := c.toLessThanHalfPi(x)

	sum, _ := evalCtx.New(xReduced)
	pow, _ := evalCtx.New(xReduced)
	den, _ := evalCtx.New(1)
	x2 := evalCtx.finalise(evalCtx.Mul(xReduced, xReduced), wpr, RoundDown, false)
	sign := int8(1)

	for i := 1; i < 150; i += 2 {
		if i > 1 {
			pow = evalCtx.finalise(evalCtx.Mul(pow, x2), wpr, RoundDown, false)
			i1, _ := evalCtx.New(int64(i - 1))
			i2, _ := evalCtx.New(int64(i))
			den = evalCtx.Mul(den, evalCtx.Mul(i1, i2))
			term := evalCtx.divide(pow, den, wpr, RoundDown, false)
			sign = -sign
			if sign < 0 {
				sum = evalCtx.Sub(sum, term)
			} else {
				sum = evalCtx.Add(sum, term)
			}
			if len(term.d) == 0 || term.d[0] == 0 {
				break
			}
		}
	}

	if quadrant > 2 {
		sum = sum.Neg()
	}

	return c.finalise(sum, c.Precision, c.Rounding, true)
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
func (c *Context) Cos(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() || x.d == nil {
		return &Decimal{s: 0}
	}
	if x.IsZero() {
		return &Decimal{s: 1, e: 0, d: []int32{1}} // cos(0) = 1
	}

	wpr := c.Precision + 20
	absE := x.e
	if absE < 0 {
		absE = -absE
	}
	if int(absE) > 0 {
		wpr += int(absE)
	}
	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	xReduced, quadrant := c.toLessThanHalfPi(x)

	sum, _ := evalCtx.New(1)
	pow, _ := evalCtx.New(1)
	den, _ := evalCtx.New(1)
	x2 := evalCtx.finalise(evalCtx.Mul(xReduced, xReduced), wpr, RoundDown, false)
	sign := int8(1)

	for i := 2; i < 150; i += 2 {
		pow = evalCtx.finalise(evalCtx.Mul(pow, x2), wpr, RoundDown, false)
		i1, _ := evalCtx.New(int64(i - 1))
		i2, _ := evalCtx.New(int64(i))
		den = evalCtx.Mul(den, evalCtx.Mul(i1, i2))
		term := evalCtx.divide(pow, den, wpr, RoundDown, false)
		sign = -sign
		if sign < 0 {
			sum = evalCtx.Sub(sum, term)
		} else {
			sum = evalCtx.Add(sum, term)
		}
		if len(term.d) == 0 || term.d[0] == 0 {
			break
		}
	}

	if quadrant == 2 || quadrant == 3 {
		sum = sum.Neg()
	}

	return c.finalise(sum, c.Precision, c.Rounding, true)
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
	if x == nil || x.IsNaN() || x.d == nil {
		return &Decimal{s: 0}
	}
	if x.IsZero() {
		return &Decimal{s: x.s, e: 0, d: []int32{0}}
	}

	wpr := c.Precision + 10
	if x.e < 0 {
		wpr += int(-x.e)
	}
	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	sinX := evalCtx.Sin(x)
	cosX := evalCtx.Cos(x)
	res := evalCtx.Div(sinX, cosX)

	return c.finalise(res, c.Precision, c.Rounding, true)
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
// Matches decimal.js inverseSine / asin (lines 909-946).
func (c *Context) Asin(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if x.IsZero() {
		return &Decimal{s: x.s, e: 0, d: []int32{0}}
	}
	if !x.IsFinite() {
		return &Decimal{s: 0}
	}

	one, _ := c.New(1)
	absX := x.Abs()
	cmpOne := absX.Cmp(one)

	if cmpOne != -1 {
		if cmpOne == 0 {
			pi := c.getPi(c.Precision + 4)
			half, _ := c.New(0.5)
			halfPi := c.Mul(pi, half)
			halfPi.s = x.s
			return c.finalise(halfPi, c.Precision, c.Rounding, true)
		}
		return &Decimal{s: 0} // NaN
	}

	wpr := c.Precision + 12
	tmpCtx := c.Clone()
	tmpCtx.Precision = wpr + 40
	tmpCtx.Rounding = RoundDown

	oneEval, _ := tmpCtx.New(1)
	oneMinusX := tmpCtx.Sub(oneEval, absX)
	if oneMinusX.e < 0 {
		wpr += int(-oneMinusX.e)
	}

	evalCtx := c.Clone()
	evalCtx.Precision = wpr + 40
	evalCtx.Rounding = RoundDown

	oneEval, _ = evalCtx.New(1)
	// asin(x) = 2 * atan(x / (sqrt((1 - x) * (1 + x)) + 1))
	oneMinusX = evalCtx.Sub(oneEval, x)
	onePlusX := evalCtx.Add(oneEval, x)
	prod := evalCtx.Mul(oneMinusX, onePlusX)
	sqrtVal := evalCtx.Sqrt(prod)
	denom := evalCtx.Add(sqrtVal, oneEval)
	frac := evalCtx.divide(x, denom, evalCtx.Precision, RoundDown, false)
	atanVal := evalCtx.Atan(frac)

	two, _ := evalCtx.New(2)
	res := evalCtx.Mul(atanVal, two)

	return c.finalise(res, c.Precision, c.Rounding, true)
}

// Asin computes inverse sine using default context.
func (x *Decimal) Asin() *Decimal {
	return globalContext.Asin(x)
}

// InverseSine is an alias for Asin.
func (x *Decimal) InverseSine() *Decimal {
	return x.Asin()
}

// Acos computes inverse cosine (acos) using context c settings.
// Matches decimal.js inverseCosine / acos (lines 725-754).
func (c *Context) Acos(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if !x.IsFinite() {
		return &Decimal{s: 0} // NaN
	}

	one, _ := c.New(1)
	absX := x.Abs()
	cmpOne := absX.Cmp(one)

	if cmpOne != -1 {
		if cmpOne == 0 {
			if x.s < 0 {
				return c.getPi(c.Precision)
			}
			return &Decimal{s: 1, e: 0, d: []int32{0}}
		}
		return &Decimal{s: 0} // NaN
	}

	if x.IsZero() {
		pi := c.getPi(c.Precision + 4)
		half, _ := c.New(0.5)
		halfPi := c.Mul(pi, half)
		return c.finalise(halfPi, c.Precision, c.Rounding, true)
	}

	wpr := c.Precision + 12
	tmpCtx := c.Clone()
	tmpCtx.Precision = wpr + 40
	tmpCtx.Rounding = RoundDown

	oneEval, _ := tmpCtx.New(1)
	oneMinusX := tmpCtx.Sub(oneEval, absX)
	if oneMinusX.e < 0 {
		wpr += int(-oneMinusX.e)
	}

	evalCtx := c.Clone()
	evalCtx.Precision = wpr + 40
	evalCtx.Rounding = RoundDown

	oneEval, _ = evalCtx.New(1)
	// acos(x) = 2 * atan(sqrt((1 - x) / (1 + x)))
	oneMinusX = evalCtx.Sub(oneEval, x)
	onePlusX := evalCtx.Add(oneEval, x)
	frac := evalCtx.divide(oneMinusX, onePlusX, evalCtx.Precision, RoundDown, false)
	sqrtVal := evalCtx.Sqrt(frac)
	atanVal := evalCtx.Atan(sqrtVal)

	two, _ := evalCtx.New(2)
	res := evalCtx.Mul(atanVal, two)

	return c.finalise(res, c.Precision, c.Rounding, true)
}

// Acos computes inverse cosine using default context.
func (x *Decimal) Acos() *Decimal {
	return globalContext.Acos(x)
}

// Atan computes inverse tangent (atan) using context c settings.
// Matches decimal.js inverseTangent / atan (lines 967-1033).
func (c *Context) Atan(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if !x.IsFinite() {
		if x.s == 0 {
			return &Decimal{s: 0}
		}
		pi := c.getPi(c.Precision + 4)
		half, _ := c.New(0.5)
		halfPi := c.Mul(pi, half)
		halfPi.s = x.s
		return c.finalise(halfPi, c.Precision, c.Rounding, true)
	}
	if x.IsZero() {
		return &Decimal{s: x.s, e: 0, d: []int32{0}}
	}

	wpr := c.Precision + 12
	wprCtx := c.Clone()
	wprCtx.Precision = wpr
	wprCtx.Rounding = RoundDown

	one, _ := wprCtx.New(1)

	// If |x| == 1
	if x.Abs().Eq(one) {
		pi := wprCtx.getPi(wpr + 4)
		quarter, _ := wprCtx.New(0.25)
		res := wprCtx.Mul(pi, quarter)
		res.s = x.s
		return c.finalise(res, c.Precision, c.Rounding, true)
	}

	// For |x| > 1: atan(x) = sign(x)*pi/2 - atan(1/|x|)
	if x.Abs().Gt(one) {
		pi := wprCtx.getPi(wpr + 4)
		half, _ := wprCtx.New(0.5)
		halfPi := wprCtx.Mul(pi, half)

		invX := wprCtx.divide(one, x.Abs(), wpr, RoundDown, false)
		atanInv := c.Atan(invX)
		res := wprCtx.Sub(halfPi, atanInv)
		res.s = x.s
		return c.finalise(res, c.Precision, c.Rounding, true)
	}

	// Argument reduction: atan(x) = 2 * atan(x / (1 + sqrt(1 + x^2)))
	k := 8
	if wprCtx.Precision > 50 {
		k = 16
	}

	xWork := x
	for i := k; i > 0; i-- {
		x2 := wprCtx.Mul(xWork, xWork)
		x2Plus1 := wprCtx.Add(x2, one)
		sqrtVal := wprCtx.Sqrt(x2Plus1)
		denom := wprCtx.Add(sqrtVal, one)
		xWork = wprCtx.divide(xWork, denom, wpr, RoundDown, false)
	}

	sum, _ := wprCtx.New(xWork)
	pow, _ := wprCtx.New(xWork)
	x2 := wprCtx.Mul(xWork, xWork)
	sign := int8(1)

	for i := 3; i < 300; i += 2 {
		pow = wprCtx.Mul(pow, x2)
		iDec, _ := wprCtx.New(int64(i))
		term := wprCtx.divide(pow, iDec, wpr, RoundDown, false)
		sign = -sign
		if sign < 0 {
			sum = wprCtx.Sub(sum, term)
		} else {
			sum = wprCtx.Add(sum, term)
		}
		if len(term.d) == 0 || term.d[0] == 0 {
			break
		}
	}

	if k > 0 {
		multVal := int64(1) << uint(k)
		multDec, _ := wprCtx.New(multVal)
		sum = wprCtx.Mul(sum, multDec)
	}

	return c.finalise(sum, c.Precision, c.Rounding, true)
}

// Atan computes inverse tangent using default context.
func (x *Decimal) Atan() *Decimal {
	return globalContext.Atan(x)
}

// Atan2 computes the arctangent of y/x using context c settings.
// Matches decimal.js atan2(y, x) (lines 4112-4153).
func (c *Context) Atan2(y, x *Decimal) *Decimal {
	if y == nil || x == nil || y.IsNaN() || x.IsNaN() {
		return &Decimal{s: 0} // NaN
	}

	pr := c.Precision
	rm := c.Rounding
	wpr := pr + 4

	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	// Both ±Infinity
	if !y.IsFinite() && !x.IsFinite() {
		pi := evalCtx.getPi(wpr)
		var factor *Decimal
		if x.s > 0 {
			factor, _ = evalCtx.New(0.25)
		} else {
			factor, _ = evalCtx.New(0.75)
		}
		r := evalCtx.Mul(pi, factor)
		r.s = y.s
		return c.finalise(r, pr, rm, false)
	}

	// x is ±Infinity or y is ±0
	if !x.IsFinite() || y.IsZero() {
		var r *Decimal
		if x.s < 0 {
			r = evalCtx.getPi(wpr)
		} else {
			r, _ = evalCtx.New(0)
		}
		r.s = y.s
		return c.finalise(r, pr, rm, false)
	}

	// y is ±Infinity or x is ±0
	if !y.IsFinite() || x.IsZero() {
		pi := evalCtx.getPi(wpr)
		half, _ := evalCtx.New(0.5)
		r := evalCtx.Mul(pi, half)
		r.s = y.s
		return c.finalise(r, pr, rm, false)
	}

	// Both non-zero and finite
	var r *Decimal
	if x.s < 0 {
		evalCtx.Precision = wpr
		evalCtx.Rounding = RoundDown

		divVal := evalCtx.divide(y, x, wpr, RoundDown, false)
		r = evalCtx.Atan(divVal)
		pi := evalCtx.getPi(wpr)

		if y.s < 0 {
			r = evalCtx.Sub(r, pi)
		} else {
			r = evalCtx.Add(r, pi)
		}
	} else {
		divVal := evalCtx.divide(y, x, wpr, RoundDown, false)
		r = evalCtx.Atan(divVal)
	}

	return c.finalise(r, pr, rm, false)
}

// Atan2 computes arctangent of y/x using default context.
func (y *Decimal) Atan2(x *Decimal) *Decimal {
	return globalContext.Atan2(y, x)
}
