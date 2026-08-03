package decimal

import (
	"math"
	"strconv"
)

func truncateDigits(d []int32, k int) ([]int32, bool) {
	if len(d) > k {
		return d[:k], true
	}
	return d, false
}

// intPow performs binary exponentiation x^n for small integer exponents n.
// Matches decimal.js intPow() (lines 3208-3241).
func (c *Context) intPow(x *Decimal, n int64, pr int) *Decimal {
	var isTruncated bool
	one, _ := c.New(1)
	r := one
	k := int(math.Ceil(float64(pr)/float64(LogBase))) + 4

	isNeg := n < 0
	if n < 0 {
		n = -n
	}

	evalCtx := c.Clone()
	evalCtx.Precision = pr + 40

	base := x

	for {
		if n%2 != 0 {
			r = evalCtx.Mul(r, base)
			var truncated bool
			r.d, truncated = truncateDigits(r.d, k)
			if truncated {
				isTruncated = true
			}
		}

		n /= 2
		if n == 0 {
			if len(r.d) > 0 {
				lastIdx := len(r.d) - 1
				if isTruncated && r.d[lastIdx] == 0 {
					r.d[lastIdx]++
				}
			}
			break
		}

		base = evalCtx.Mul(base, base)
		base.d, _ = truncateDigits(base.d, k)
	}

	if isNeg {
		return c.Div(one, r)
	}

	return r
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

	// pow(x, 1) = x
	one, _ := c.New(1)
	if y.Eq(one) {
		return c.finalise(new(Decimal).Set(x), c.Precision, c.Rounding, false)
	}

	// pow(x, 0.5) = sqrt(x)
	half, _ := c.New(0.5)
	if y.Eq(half) {
		return c.Sqrt(x)
	}

	// Non-finite or zero operands: follow math.Pow(float64(x), float64(y))
	if x.d == nil || y.d == nil || x.IsZero() {
		fx, errX := x.Float64()
		fy, errY := y.Float64()
		if errX == nil && errY == nil {
			resF := math.Pow(fx, fy)
			if math.IsNaN(resF) {
				return &Decimal{s: 0}
			}
			if math.IsInf(resF, 0) {
				s := int8(1)
				if resF < 0 {
					s = -1
				}
				return &Decimal{s: s, e: 0, d: nil}
			}
			if resF == 0 {
				s := int8(1)
				if math.Signbit(resF) {
					s = -1
				}
				return &Decimal{s: s, e: 0, d: []int32{0}}
			}
			resDec, err := c.New(strconv.FormatFloat(resF, 'g', -1, 64))
			if err == nil {
				return resDec
			}
		}
	}

	if x.IsNeg() {
		if !y.IsInt() {
			return &Decimal{s: 0}
		}
		xAbs := &Decimal{s: 1, e: x.e, d: x.d}
		res := c.Pow(xAbs, y)
		if res == nil || res.IsNaN() {
			return res
		}
		// Determine parity of y
		var isOdd bool
		if len(y.d) > 0 {
			lastLimb := y.d[len(y.d)-1]
			isOdd = lastLimb%2 != 0
		}
		if isOdd {
			res.s = -1
		} else {
			res.s = 1
		}
		return res
	}

	// Check if y is an integer or half-integer with exact square base
	two, _ := c.New(2)
	twoY := c.Mul(y, two)
	if twoY.IsInt() {
		if n, err := strconv.ParseInt(twoY.String(), 10, 64); err == nil {
			if n >= -MaxSafeInteger && n <= MaxSafeInteger {
				if y.IsInt() {
					r := c.intPow(x, n/2, c.Precision)
					return c.finalise(r, c.Precision, c.Rounding, false)
				}
				sqrtX := c.Sqrt(x)
				if c.Mul(sqrtX, sqrtX).Eq(x) {
					r := c.intPow(sqrtX, n, c.Precision)
					return c.finalise(r, c.Precision, c.Rounding, false)
				}
			}
		}
	}

	// x^y = exp(y * ln(x))
	evalCtx := c.Clone()
	evalCtx.Precision = c.Precision + len(x.String()) + len(y.String()) + 100
	lnX := evalCtx.Ln(x)
	yLnX := evalCtx.Mul(y, lnX)
	resExp := evalCtx.Exp(yLnX)
	return c.finalise(resExp, c.Precision, c.Rounding, true)
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
