package decimal

import (
	"math"
)

// divide performs division in base 10^7 with precision pr and rounding mode rm.
// Matches decimal.js divide() (lines 2728-2939).
func (c *Context) divide(x, y *Decimal, pr int, rm RoundingMode, dp bool) *Decimal {
	sign := x.s * y.s

	// Non-finite handling
	if x.d == nil || y.d == nil || len(x.d) == 0 || len(y.d) == 0 || x.d[0] == 0 || y.d[0] == 0 {
		if x.s == 0 || y.s == 0 || (x.d == nil && y.d == nil) || (len(x.d) > 0 && x.d[0] == 0 && len(y.d) > 0 && y.d[0] == 0) {
			return &Decimal{s: 0}
		}
		if (len(x.d) > 0 && x.d[0] == 0) || y.d == nil {
			return &Decimal{s: sign, e: 0, d: []int32{0}}
		}
		return &Decimal{s: sign, e: 0, d: nil}
	}

	eX := int64(math.Floor(float64(x.e) / float64(LogBase)))
	eY := int64(math.Floor(float64(y.e) / float64(LogBase)))
	e := eX - eY

	xd := append([]int32(nil), x.d...)
	yd := append([]int32(nil), y.d...)
	xL := len(xd)
	yL := len(yd)

	// Adjust exponent estimate if divisor magnitude is larger than dividend
	i := 0
	for i < yL && i < xL && yd[i] == xd[i] {
		i++
	}
	if i < yL && yd[i] > func() int32 {
		if i < xL {
			return xd[i]
		}
		return 0
	}() {
		e--
	}

	sd := pr
	if dp {
		sd = pr + int(x.e-y.e) + 1
	}

	if sd < 0 {
		return c.finalise(&Decimal{s: sign, e: 0, d: []int32{1}}, pr, rm, true)
	}

	// Calculate required limbs
	sdLimbs := sd/LogBase + 2
	var qd []int32
	var isTruncated bool

	if yL == 1 {
		// Single-limb divisor fast path
		ydScalar := yd[0]
		sdLimbs++
		var k int64 = 0
		for idx := 0; (idx < xL || k > 0) && sdLimbs > 0; idx++ {
			var xiVal int64 = 0
			if idx < xL {
				xiVal = int64(xd[idx])
			}
			t := k*int64(Base) + xiVal
			qd = append(qd, int32(t/int64(ydScalar)))
			k = t % int64(ydScalar)
			sdLimbs--
		}
		isTruncated = k > 0
	} else {
		// Normalization factor k to ensure yd[0] >= Base / 2
		normK := Base / (yd[0] + 1)
		if normK > 1 {
			yd = multiplyInteger(yd, normK, Base)
			xd = multiplyInteger(xd, normK, Base)
			yL = len(yd)
			xL = len(xd)
		}

		rem := append([]int32(nil), xd...)
		if len(rem) < yL {
			zeros := make([]int32, yL-len(rem))
			rem = append(rem, zeros...)
		}

		xi := yL
		yd0 := yd[0]
		if yd[1] >= Base/2 {
			yd0++
		}

		for sdLimbs >= 0 {
			remL := len(rem)
			cmp := compareLimbs(yd, rem, yL, remL)

			var k int32 = 0
			if cmp < 0 {
				rem0 := int64(rem[0])
				if yL != remL && len(rem) > 1 {
					rem0 = rem0*int64(Base) + int64(rem[1])
				}
				k = int32(rem0 / int64(yd0))

				if k > 1 {
					if k >= Base {
						k = Base - 1
					}
					prod := multiplyInteger(yd, k, Base)
					prodL := len(prod)
					cmpProd := compareLimbs(prod, rem, prodL, remL)
					if cmpProd == 1 {
						k--
					}
				} else {
					k = 1
				}
			} else if cmp == 0 {
				k = 1
				rem = []int32{0}
			}

			qd = append(qd, k)
			if xi < xL {
				rem = append(rem, xd[xi])
				xi++
			} else {
				rem = append(rem, 0)
			}
			sdLimbs--
		}

		isTruncated = len(rem) > 0 && rem[0] != 0
	}

	// Leading zero shift
	if len(qd) > 0 && qd[0] == 0 {
		qd = qd[1:]
	}

	// Calculate base 10 exponent
	var qExp int64 = 0
	if len(qd) > 0 {
		iDigits := 1
		for kVal := qd[0]; kVal >= 10; kVal /= 10 {
			iDigits++
		}
		qExp = int64(iDigits) + e*int64(LogBase) - 1
	}

	res := &Decimal{
		s: sign,
		e: qExp,
		d: qd,
	}

	targetPr := pr
	if dp {
		targetPr = pr + int(res.e) + 1
	}

	return c.finalise(res, targetPr, rm, isTruncated)
}

// Div computes x / y using context c settings.
// Matches decimal.js dividedBy / div (lines 468-470).
func (c *Context) Div(x, y *Decimal) *Decimal {
	return c.divide(x, y, c.Precision, c.Rounding, false)
}

// Div computes x / y using default context.
func (x *Decimal) Div(y *Decimal) *Decimal {
	return globalContext.Div(x, y)
}

// DividedBy is an alias for Div.
func (x *Decimal) DividedBy(y *Decimal) *Decimal {
	return x.Div(y)
}
