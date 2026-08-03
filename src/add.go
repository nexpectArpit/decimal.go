package decimal

import (
	"math"
)

// Add computes x + y using context c settings.
// Matches decimal.js P.plus / P.add (lines 1262-1358).
func (c *Context) Add(x, y *Decimal) *Decimal {
	if x == nil || y == nil {
		return &Decimal{s: 0}
	}

	xVal := x
	yVal := y

	// If either is not finite...
	if xVal.d == nil || yVal.d == nil {
		if xVal.s == 0 || yVal.s == 0 {
			return &Decimal{s: 0}
		}
		if xVal.d == nil && yVal.d == nil {
			if xVal.s == yVal.s {
				return &Decimal{s: xVal.s, e: 0, d: nil}
			}
			return &Decimal{s: 0} // Inf - Inf = NaN
		}
		if xVal.d == nil {
			return &Decimal{s: xVal.s, e: 0, d: nil}
		}
		return &Decimal{s: yVal.s, e: 0, d: nil}
	}

	// If signs differ, perform subtraction: x + (-y) -> x - y
	if xVal.s != yVal.s {
		negY := &Decimal{s: -yVal.s, e: yVal.e, d: yVal.d}
		return c.Sub(xVal, negY)
	}

	// If either is zero...
	if (len(xVal.d) == 0 || xVal.d[0] == 0) && (len(yVal.d) == 0 || yVal.d[0] == 0) {
		resSign := int8(1)
		if xVal.s < 0 && yVal.s < 0 {
			resSign = -1
		} else if c.Rounding == RoundFloor {
			resSign = -1
		}
		return &Decimal{s: resSign, e: 0, d: []int32{0}}
	}
	if len(xVal.d) == 0 || xVal.d[0] == 0 {
		res := &Decimal{s: yVal.s, e: yVal.e, d: append([]int32(nil), yVal.d...)}
		return c.finalise(res, c.Precision, c.Rounding, false)
	}
	if len(yVal.d) == 0 || yVal.d[0] == 0 {
		res := &Decimal{s: xVal.s, e: xVal.e, d: append([]int32(nil), xVal.d...)}
		return c.finalise(res, c.Precision, c.Rounding, false)
	}

	// Exponents in base 10^7
	eY := int64(math.Floor(float64(yVal.e) / float64(LogBase)))
	eX := int64(math.Floor(float64(xVal.e) / float64(LogBase)))

	xd := append([]int32(nil), xVal.d...)
	yd := append([]int32(nil), yVal.d...)
	k := eX - eY

	e := eY
	if k != 0 {
		var d *[]int32
		if k < 0 {
			lenOther := len(yd)
			limit := int64(math.Max(math.Ceil(float64(c.Precision)/float64(LogBase)), float64(lenOther))) + 1
			if -k > limit {
				return c.finalise(yVal, c.Precision, c.Rounding, false)
			}
			d = &xd
			k = -k
		} else {
			lenOther := len(xd)
			limit := int64(math.Max(math.Ceil(float64(c.Precision)/float64(LogBase)), float64(lenOther))) + 1
			if k > limit {
				return c.finalise(xVal, c.Precision, c.Rounding, false)
			}
			d = &yd
			e = eX
		}

		// Prepend zeros to align exponents
		zeros := make([]int32, k)
		*d = append(zeros, *d...)
	}

	// Make limb slices same length for addition
	if len(xd) < len(yd) {
		zeros := make([]int32, len(yd)-len(xd))
		xd = append(xd, zeros...)
	} else if len(yd) < len(xd) {
		zeros := make([]int32, len(xd)-len(yd))
		yd = append(yd, zeros...)
	}

	// Add limb by limb from right to left
	var carry int32 = 0
	for i := len(xd) - 1; i >= 0; i-- {
		sum := xd[i] + yd[i] + carry
		if sum >= Base {
			sum -= Base
			carry = 1
		} else {
			carry = 0
		}
		xd[i] = sum
	}

	if carry > 0 {
		xd = append([]int32{carry}, xd...)
		e++
	}

	// Remove trailing zero limbs
	for last := len(xd) - 1; last >= 0 && xd[last] == 0; last-- {
		xd = xd[:last]
	}

	res := &Decimal{
		s: xVal.s,
		e: getBase10Exponent(xd, e),
		d: xd,
	}

	return c.finalise(res, c.Precision, c.Rounding, false)
}

// Add computes x + y using default context.
func (x *Decimal) Add(y *Decimal) *Decimal {
	return globalContext.Add(x, y)
}

// Plus is an alias for Add.
func (x *Decimal) Plus(y *Decimal) *Decimal {
	return x.Add(y)
}
