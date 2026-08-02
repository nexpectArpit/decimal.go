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

	xVal := x
	yVal := y

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
		negY := &Decimal{s: -yVal.s, e: yVal.e, d: yVal.d}
		return c.Add(xVal, negY)
	}

	// Zero handling
	xZero := len(xVal.d) > 0 && xVal.d[0] == 0
	yZero := len(yVal.d) > 0 && yVal.d[0] == 0

	if xZero || yZero {
		if yZero && !xZero {
			return c.finalise(xVal, c.Precision, c.Rounding, false)
		}
		if xZero && !yZero {
			negY := &Decimal{s: -yVal.s, e: yVal.e, d: yVal.d}
			return c.finalise(negY, c.Precision, c.Rounding, false)
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

		var lenOther int
		if xLTy {
			lenOther = len(yd)
		} else {
			lenOther = len(xd)
		}
		limit := int64(math.Max(math.Ceil(float64(c.Precision)/float64(LogBase)), float64(lenOther))) + 2
		if k > limit {
			if xLTy {
				negY := &Decimal{
					s: -yVal.s,
					e: yVal.e,
					d: yVal.d,
				}
				return c.finalise(negY, c.Precision, c.Rounding, false)
			}
			return c.finalise(xVal, c.Precision, c.Rounding, false)
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

	resSign := xVal.s
	if xLTy {
		xd, yd = yd, xd
		resSign = -xVal.s
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
		s: resSign,
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
// Matches decimal.js abs() (lines 530-540).
func (x *Decimal) Abs() *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if x.d == nil {
		return &Decimal{s: 1, e: 0, d: nil}
	}
	dCopy := make([]int32, len(x.d))
	copy(dCopy, x.d)
	return &Decimal{s: 1, e: x.e, d: dCopy}
}

// Neg returns a new Decimal representing the negated value of x (-x).
func (x *Decimal) Neg() *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	var dCopy []int32
	if x.d != nil {
		dCopy = make([]int32, len(x.d))
		copy(dCopy, x.d)
	}
	s := x.s
	if s != 0 {
		s = -s
	}
	return &Decimal{s: s, e: x.e, d: dCopy}
}
