package decimal

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func initialSqrtEstimate(c *Context, x *Decimal) *Decimal {
	// Try float64 conversion
	f, err := x.Float64()
	if err == nil && f > 0 && !math.IsInf(f, 0) {
		rVal, err2 := c.New(strconv.FormatFloat(math.Sqrt(f), 'g', -1, 64))
		if err2 == nil {
			return rVal
		}
	}

	// Underflow or overflow case
	n := digitsToString(x.d)
	e := x.e

	if (int64(len(n))+e)%2 == 0 {
		n += "0"
	}

	fVal, _ := strconv.ParseFloat(n, 64)
	s := math.Sqrt(fVal)

	var newE int64
	if e >= 0 {
		newE = (e + 1) / 2
	} else {
		newE = (e + 1 - 1) / 2
	}
	if e < 0 || e%2 != 0 {
		newE--
	}

	var nStr string
	if math.IsInf(s, 0) {
		nStr = fmt.Sprintf("5e%d", newE)
	} else {
		expStr := strconv.FormatFloat(s, 'e', -1, 64)
		eIdx := strings.IndexByte(expStr, 'e')
		if eIdx > 0 {
			nStr = expStr[:eIdx+1] + strconv.FormatInt(newE, 10)
		} else {
			nStr = expStr
		}
	}

	rVal, _ := c.New(nStr)
	return rVal
}

// Sqrt computes square root sqrt(x) using context c settings.
// Matches decimal.js squareRoot / sqrt (lines 1432-1522).
func (c *Context) Sqrt(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if x.IsZero() || x.d == nil {
		return &Decimal{s: x.s, e: x.e, d: x.d}
	}
	if x.s < 0 {
		return &Decimal{s: 0} // sqrt(-positive) = NaN
	}

	evalCtx := c.Clone()
	evalCtx.Precision = c.Precision + 16
	sd := c.Precision + 3
	half, _ := evalCtx.New(0.5)

	r := initialSqrtEstimate(evalCtx, x)
	if r == nil || r.IsZero() {
		r, _ = evalCtx.New(x)
	}

	var m bool
	var rep bool

	for i := 0; i < 60; i++ {
		t := r
		xDivR := evalCtx.divide(x, t, sd+2, RoundDown, false)
		sum := evalCtx.Add(t, xDivR)
		r = evalCtx.Mul(half, sum)

		tDigits := digitsToString(t.d)
		rDigits := digitsToString(r.d)

		if len(tDigits) >= sd && len(rDigits) >= sd && tDigits[:sd] == rDigits[:sd] {
			var n string
			if len(rDigits) >= sd+1 {
				start := sd - 3
				if start < 0 {
					start = 0
				}
				n = rDigits[start : sd+1]
			}

			if n == "9999" || (!rep && n == "4999") {
				if !rep {
					tCheck := c.finalise(t, c.Precision+1, RoundDown, false)
					if c.Mul(tCheck, tCheck).Eq(x) {
						r = t
						break
					}
				}
				sd += 4
				rep = true
			} else {
				if n == "0000" || (len(n) > 1 && n[0] == '5' && n[1:] == "000") {
					rCheck := c.finalise(r, c.Precision+1, RoundDown, false)
					m = !c.Mul(rCheck, rCheck).Eq(x)
				}
				break
			}
		}
	}

	return c.finalise(r, c.Precision, c.Rounding, m)
}

// Sqrt computes square root sqrt(x) using default context.
func (x *Decimal) Sqrt() *Decimal {
	return globalContext.Sqrt(x)
}

// SquareRoot is an alias for Sqrt.
func (x *Decimal) SquareRoot() *Decimal {
	return x.Sqrt()
}

func initialCbrtEstimate(c *Context, x *Decimal) *Decimal {
	f, err := x.Float64()
	if err == nil && f != 0 && !math.IsInf(f, 0) {
		cbrtVal := math.Cbrt(f)
		rVal, err2 := c.New(strconv.FormatFloat(cbrtVal, 'g', -1, 64))
		if err2 == nil {
			return rVal
		}
	}

	n := digitsToString(x.d)
	e := x.e

	mod3 := (e - int64(len(n)) + 1) % 3
	if mod3 == 1 || mod3 == -2 {
		n += "0"
	} else if mod3 != 0 {
		n += "00"
	}

	fVal, _ := strconv.ParseFloat(n, 64)
	s := math.Cbrt(fVal)

	newE := int64(math.Floor(float64(e+1) / 3.0))
	if e%3 == func() int64 {
		if e < 0 {
			return -1
		}
		return 2
	}() {
		newE--
	}

	var nStr string
	if math.IsInf(s, 0) || s == 0 {
		nStr = fmt.Sprintf("5e%d", newE)
	} else {
		expStr := strconv.FormatFloat(s, 'e', -1, 64)
		eIdx := strings.IndexByte(expStr, 'e')
		if eIdx > 0 {
			nStr = expStr[:eIdx+1] + strconv.FormatInt(newE, 10)
		} else {
			nStr = expStr
		}
	}

	rVal, _ := c.New(nStr)
	if rVal != nil && x.s < 0 {
		rVal.s = -1
	}
	return rVal
}

// Cbrt computes cube root cbrt(x) using context c settings.
// Matches decimal.js cubeRoot / cbrt (lines 334-421).
func (c *Context) Cbrt(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if x.IsZero() || x.d == nil {
		return &Decimal{s: x.s, e: x.e, d: x.d}
	}

	wpr := c.Precision + 12
	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	r := initialCbrtEstimate(evalCtx, x)
	if r == nil || r.IsZero() {
		r, _ = evalCtx.New(x)
	}
	two, _ := evalCtx.New(2)

	// Halley's cubic method: r_{k+1} = r_k * (r_k^3 + 2X) / (2r_k^3 + X)
	for i := 0; i < 60; i++ {
		r2 := evalCtx.Mul(r, r)
		r3 := evalCtx.Mul(r2, r)
		twoX := evalCtx.Mul(two, x)
		num := evalCtx.Add(r3, twoX)

		twoR3 := evalCtx.Mul(two, r3)
		den := evalCtx.Add(twoR3, x)

		frac := evalCtx.divide(num, den, wpr, RoundDown, false)
		nextR := evalCtx.Mul(r, frac)
		if nextR.Eq(r) {
			break
		}
		r = nextR
	}

	return c.finalise(r, c.Precision, c.Rounding, true)
}

// Cbrt computes cube root cbrt(x) using default context.
func (x *Decimal) Cbrt() *Decimal {
	return globalContext.Cbrt(x)
}

// CubeRoot is an alias for Cbrt.
func (x *Decimal) CubeRoot() *Decimal {
	return x.Cbrt()
}

// Hypot computes the square root of the sum of squares of its arguments.
// Matches decimal.js hypot (lines 4522-4544).
func (c *Context) Hypot(args ...*Decimal) *Decimal {
	if len(args) == 0 {
		return &Decimal{s: 1, e: 0, d: []int32{0}}
	}

	evalCtx := c.Clone()
	evalCtx.Precision = c.Precision + 16

	zero, _ := evalCtx.New(0)
	sum := zero

	var hasInf bool
	var hasNaN bool

	for _, arg := range args {
		if arg == nil || arg.IsNaN() {
			hasNaN = true
			continue
		}
		if arg.d == nil {
			hasInf = true
			continue
		}
		sq := evalCtx.Mul(arg, arg)
		sum = evalCtx.Add(sum, sq)
	}

	if hasInf {
		return &Decimal{s: 1, e: 0, d: nil}
	}
	if hasNaN {
		return &Decimal{s: 0}
	}

	res := evalCtx.Sqrt(sum)
	return c.finalise(res, c.Precision, c.Rounding, false)
}
