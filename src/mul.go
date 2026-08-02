package decimal

import (
	"math"
)

// multiplyInteger multiplies limb slice x by integer scalar k in base.
// Matches decimal.js multiplyInteger() (lines 2681-2695).
func multiplyInteger(x []int32, k int32, base int32) []int32 {
	if len(x) == 0 || k == 0 {
		return []int32{0}
	}

	res := make([]int32, len(x))
	copy(res, x)
	var carry int64 = 0
	b64 := int64(base)
	k64 := int64(k)

	for i := len(res) - 1; i >= 0; i-- {
		temp := int64(res[i])*k64 + carry
		res[i] = int32(temp % b64)
		carry = temp / b64
	}

	if carry > 0 {
		res = append([]int32{int32(carry)}, res...)
	}

	return res
}

// Mul computes x * y using context c settings.
// Matches decimal.js times / mul (lines 1577-1678).
func (c *Context) Mul(x, y *Decimal) *Decimal {
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
		if (xVal.d == nil && len(yVal.d) > 0 && yVal.d[0] == 0) ||
			(yVal.d == nil && len(xVal.d) > 0 && xVal.d[0] == 0) {
			return &Decimal{s: 0} // Inf * 0 = NaN
		}
		sign := xVal.s * yVal.s
		return &Decimal{s: sign, e: 0, d: nil}
	}

	sign := xVal.s * yVal.s

	// Zero handling
	if (len(xVal.d) > 0 && xVal.d[0] == 0) || (len(yVal.d) > 0 && yVal.d[0] == 0) {
		return &Decimal{s: sign, e: 0, d: []int32{0}}
	}

	eX := int64(math.Floor(float64(xVal.e) / float64(LogBase)))
	eY := int64(math.Floor(float64(yVal.e) / float64(LogBase)))

	xd := xVal.d
	yd := yVal.d
	xL := len(xd)
	yL := len(yd)

	// Result exponent estimation
	e := eX + eY

	// $O(N \cdot M)$ Schoolbook limb multiplication
	resL := xL + yL
	resLimbs := make([]int64, resL)

	for i := xL - 1; i >= 0; i-- {
		for j := yL - 1; j >= 0; j-- {
			resLimbs[i+j+1] += int64(xd[i]) * int64(yd[j])
		}
	}

	// Carry propagation in base 10^7
	b64 := int64(Base)
	for i := resL - 1; i > 0; i-- {
		resLimbs[i-1] += resLimbs[i] / b64
		resLimbs[i] %= b64
	}

	if resLimbs[0] == 0 {
		resLimbs = resLimbs[1:]
	} else {
		e++
	}

	// Convert back to int32 slice
	finalLimbs := make([]int32, len(resLimbs))
	for i, v := range resLimbs {
		finalLimbs[i] = int32(v)
	}

	// Remove trailing zero limbs
	for last := len(finalLimbs) - 1; last >= 0 && finalLimbs[last] == 0; last-- {
		finalLimbs = finalLimbs[:last]
	}

	res := &Decimal{
		s: sign,
		e: getBase10Exponent(finalLimbs, e),
		d: finalLimbs,
	}

	return c.finalise(res, c.Precision, c.Rounding, false)
}

// Mul computes x * y using default context.
func (x *Decimal) Mul(y *Decimal) *Decimal {
	return globalContext.Mul(x, y)
}

// Times is an alias for Mul.
func (x *Decimal) Times(y *Decimal) *Decimal {
	return x.Mul(y)
}
