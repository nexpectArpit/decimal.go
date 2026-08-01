package decimal

import (
	"math"
)

// Sub computes x - y using context c settings.
// Matches decimal.js P.minus / P.sub (lines 1264-1411).
func (c *Context) Sub(x, y *Decimal) *Decimal {
	if x == nil || y == nil {
		return &Decimal{s: 0}
	}

	xVal, _ := c.New(x)
	yVal, _ := c.New(y)

	// Non-finite handling
	if xVal.d == nil || yVal.d == nil {
		if xVal.s == 0 || yVal.s == 0 {
			return &Decimal{s: 0}
		}
		if xVal.d == nil && yVal.d == nil {
			if xVal.s != yVal.s {
				return &Decimal{s: xVal.s, e: 0, d: nil}
			}
			return &Decimal{s: 0} // Inf - Inf = NaN
		}
		if xVal.d == nil {
			return &Decimal{s: xVal.s, e: 0, d: nil}
		}
		return &Decimal{s: -yVal.s, e: 0, d: nil}
	}

	// If signs differ, convert to addition: x - (-y) -> x + y
	if xVal.s != yVal.s {
		yVal.s = -yVal.s
		return c.Add(xVal, yVal)
	}

	// Zero handling
	xZero := len(xVal.d) > 0 && xVal.d[0] == 0
	yZero := len(yVal.d) > 0 && yVal.d[0] == 0

	if xZero || yZero {
		if yZero && !xZero {
			return c.finalise(xVal, c.Precision, c.Rounding, false)
		}
		if xZero && !yZero {
			yVal.s = -yVal.s
			return c.finalise(yVal, c.Precision, c.Rounding, false)
		}
		s := int8(1)
		if c.Rounding == RoundFloor {
			s = -1
		}
		return &Decimal{s: s, e: 0, d: []int32{0}}
	}

	eY := int64(math.Floor(float64(yVal.e) / float64(LogBase)))
	eX := int64(math.Floor(float64(xVal.e) / float64(LogBase)))

	xd := append([]int32(nil), xVal.d...)
	yd := append([]int32(nil), yVal.d...)
	k := eX - eY

	e := eY
	xLTy := false

	if k != 0 {
		xLTy = k < 0
		var d *[]int32
		if xLTy {
			d = &xd
			k = -k
		} else {
			d = &yd
			e = eX
		}

		limit := int64(math.Ceil(float64(c.Precision)/float64(LogBase))) + 2
		if k > limit {
			k = limit
			*d = (*d)[:1]
		}

		zeros := make([]int32, k)
		*d = append(zeros, *d...)
	} else {
		// Compare digit slices to find larger magnitude
		cmpRes := compareLimbs(xd, yd, len(xd), len(yd))
		if cmpRes < 0 {
			xLTy = true
		}
	}

	if xLTy {
		xd, yd = yd, xd
		yVal.s = -yVal.s
	}

	// Append zeros to xd if shorter than yd
	if len(xd) < len(yd) {
		zeros := make([]int32, len(yd)-len(xd))
		xd = append(xd, zeros...)
	}

	// Subtract yd from xd
	for i := len(yd) - 1; i >= 0; i-- {
		if xd[i] < yd[i] {
			j := i
			for j > 0 && xd[j-1] == 0 {
				j--
				xd[j] = Base - 1
			}
			if j > 0 {
				xd[j-1]--
			}
			xd[i] += Base
		}
		xd[i] -= yd[i]
	}

	// Remove trailing zero limbs
	for last := len(xd) - 1; last >= 0 && xd[last] == 0; last-- {
		xd = xd[:last]
	}

	// Remove leading zero limbs and adjust exponent
	for len(xd) > 0 && xd[0] == 0 {
		xd = xd[1:]
		e--
	}

	if len(xd) == 0 || xd[0] == 0 {
		s := int8(1)
		if c.Rounding == RoundFloor {
			s = -1
		}
		return &Decimal{s: s, e: 0, d: []int32{0}}
	}

	res := &Decimal{
		s: yVal.s,
		e: getBase10Exponent(xd, e),
		d: xd,
	}

	return c.finalise(res, c.Precision, c.Rounding, false)
}

// Sub computes x - y using default context.
func (x *Decimal) Sub(y *Decimal) *Decimal {
	return globalContext.Sub(x, y)
}

// Minus is an alias for Sub.
func (x *Decimal) Minus(y *Decimal) *Decimal {
	return x.Sub(y)
}

// Abs returns a new Decimal representing the absolute value of x.
func (x *Decimal) Abs() *Decimal {
	if x == nil {
		return &Decimal{s: 0}
	}
	res, _ := globalContext.New(x)
	if res.s < 0 {
		res.s = 1
	}
	return globalContext.finalise(res, globalContext.Precision, globalContext.Rounding, false)
}

// Neg returns a new Decimal representing the negated value of x (-x).
func (x *Decimal) Neg() *Decimal {
	if x == nil {
		return &Decimal{s: 0}
	}
	res, _ := globalContext.New(x)
	if res.s != 0 {
		res.s = -res.s
	}
	return globalContext.finalise(res, globalContext.Precision, globalContext.Rounding, false)
}
