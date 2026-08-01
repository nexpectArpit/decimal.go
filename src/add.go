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

	// Make copies of inputs
	xVal, _ := c.New(x)
	yVal, _ := c.New(y)

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
		yVal.s = -yVal.s
		return c.Sub(xVal, yVal)
	}

	// If either is zero...
	if len(xVal.d) > 0 && xVal.d[0] == 0 {
		if len(yVal.d) > 0 && yVal.d[0] == 0 {
			if c.Rounding == RoundFloor {
				yVal.s = -1
			} else {
				yVal.s = 1
			}
			return c.finalise(yVal, c.Precision, c.Rounding, false)
		}
		return c.finalise(yVal, c.Precision, c.Rounding, false)
	}
	if len(yVal.d) > 0 && yVal.d[0] == 0 {
		return c.finalise(xVal, c.Precision, c.Rounding, false)
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
