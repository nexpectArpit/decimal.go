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

	// If both are zero (same sign, since opposite signs were redirected to Sub above)...
	// x + x retains the sign of x even when x is zero (IEEE 754-2008 section 6.3);
	// the ROUND_FLOOR sign rule for an exact-zero sum only applies to opposite-signed
	// operands, which is handled by the Sub() path above, not here.
	if (len(xVal.d) == 0 || xVal.d[0] == 0) && (len(yVal.d) == 0 || yVal.d[0] == 0) {
		return &Decimal{s: xVal.s, e: 0, d: []int32{0}}
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
		var lenOther int
		if k < 0 {
			lenOther = len(yd)
			d = &xd
			k = -k
		} else {
			lenOther = len(xd)
			d = &yd
			e = eX
		}

		// Limit the number of zeros prepended to max(ceil(pr/LOG_BASE), lenOther) + 1,
		// truncating the smaller operand to its leading limb rather than dropping it
		// entirely, so an astronomically smaller operand still contributes a tiny
		// nonzero remainder that can trigger correct borrow/rounding behavior
		// (matches decimal.js P.plus, lines 1591-1603).
		limit := int64(math.Ceil(float64(c.Precision) / float64(LogBase)))
		if limit <= int64(lenOther) {
			limit = int64(lenOther) + 1
		} else {
			limit++
		}
		if k > limit {
			k = limit
			if len(*d) > 1 {
				*d = (*d)[:1]
			}
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
