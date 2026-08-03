package decimal

import (
	"math"
)

// finalise rounds Decimal x to sd significant digits using rounding mode rm.
// Matches decimal.js finalise() (lines 2946-3110).
func (c *Context) finalise(x *Decimal, sd int, rm RoundingMode, isTruncated bool) *Decimal {
	xd := x.d
	if xd == nil {
		return x
	}

	var digits, i, j, k, rd int
	var roundUp bool
	var w int32
	var xdi int

	// Determine length of first word xd[0]
	digits = 1
	for kVal := xd[0]; kVal >= 10; kVal /= 10 {
		digits++
	}
	i = sd - digits

	// Rounding digit in first word?
	if i < 0 {
		i += LogBase
		j = sd
		xdi = 0
		w = xd[0]
		rd = int(w/int32(math.Pow10(digits-j-1))) % 10
	} else {
		xdi = int(math.Ceil(float64(i+1) / float64(LogBase)))
		kLength := len(xd)
		if xdi >= kLength {
			if isTruncated {
				for kLength <= xdi {
					xd = append(xd, 0)
					kLength++
				}
				x.d = xd
				w = 0
				rd = 0
				digits = 1
				i %= LogBase
				j = i - LogBase + 1
			} else {
				goto checkBounds
			}
		} else {
			w = xd[xdi]
			kVal := w
			digits = 1
			for kVal >= 10 {
				kVal /= 10
				digits++
			}
			i %= LogBase
			j = i - LogBase + digits
			if j < 0 {
				rd = 0
			} else {
				rd = int(w/int32(math.Pow10(digits-j-1))) % 10
			}
		}
	}

	// Any non-zero digits after rounding digit?
	if !isTruncated {
		if sd < 0 || xdi+1 < len(xd) {
			isTruncated = true
		} else if j < 0 {
			isTruncated = w != 0
		} else {
			rem := w % int32(math.Pow10(digits-j-1))
			isTruncated = rem != 0
		}
	}

	if rm < RoundHalfUp {
		targetMode := RoundCeil
		if x.s < 0 {
			targetMode = RoundFloor
		}
		roundUp = (rd != 0 || isTruncated) && (rm == RoundUp || rm == targetMode)
	} else {
		var leftDigit int
		if i > 0 {
			if j > 0 {
				leftDigit = int(w / int32(math.Pow10(digits-j)))
			} else {
				leftDigit = 0
			}
		} else if xdi > 0 {
			leftDigit = int(xd[xdi-1]) % 10
		}

		isOddLeft := leftDigit%2 != 0
		targetHalfMode := RoundHalfCeil
		if x.s < 0 {
			targetHalfMode = RoundHalfFloor
		}

		roundUp = rd > 5 || (rd == 5 && (rm == RoundHalfUp || isTruncated || (rm == RoundHalfEven && isOddLeft) || rm == targetHalfMode))
	}

	if sd < 1 || len(xd) == 0 || xd[0] == 0 {
		if roundUp {
			newSd := int64(sd) - (x.e + 1)
			rem := newSd % int64(LogBase)
			if rem < 0 {
				rem += int64(LogBase)
			}
			exp := (int64(LogBase) - rem) % int64(LogBase)
			x.d = []int32{int32(math.Pow10(int(exp)))}
			x.e = -newSd
		} else {
			x.d = []int32{0}
			x.e = 0
		}
		return x
	}

	// Remove excess digits
	if i == 0 {
		x.d = xd[:xdi]
		k = 1
		xdi--
	} else {
		x.d = xd[:xdi+1]
		k = int(math.Pow10(LogBase - i))
		if j > 0 {
			wPart := int32(float64(w) / math.Pow10(digits-j))
			wPart %= int32(math.Pow10(j))
			x.d[xdi] = wPart * int32(k)
		} else {
			x.d[xdi] = 0
		}
	}

	if roundUp {
		for {
			if xdi == 0 {
				iLen := 1
				for jVal := x.d[0]; jVal >= 10; jVal /= 10 {
					iLen++
				}
				x.d[0] += int32(k)
				kLen := 1
				for jVal := x.d[0]; jVal >= 10; jVal /= 10 {
					kLen++
				}
				if iLen != kLen {
					x.e++
					if x.d[0] == Base {
						x.d[0] = 1
					}
				}
				break
			} else {
				x.d[xdi] += int32(k)
				if x.d[xdi] < Base {
					break
				}
				carry := x.d[xdi] / Base
				x.d[xdi] %= Base
				xdi--
				k = int(carry)
			}
		}
	}

	// Remove trailing zero limbs only if all limbs are zero
	if len(x.d) > 0 {
		allZero := true
		for _, v := range x.d {
			if v != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			x.d = []int32{0}
		} else {
			for len(x.d) > 1 && x.d[len(x.d)-1] == 0 {
				x.d = x.d[:len(x.d)-1]
			}
		}
	}

checkBounds:
	if x.e > c.MaxE {
		x.d = nil
		x.e = 0
	} else if x.e < c.MinE {
		x.e = 0
		x.d = []int32{0}
	}

	return x
}

// Trunc rounds Decimal x towards zero to an integer.
func (c *Context) Trunc(x *Decimal) *Decimal {
	if x == nil || !x.IsFinite() || x.IsZero() {
		return x
	}
	if x.e < 0 {
		return &Decimal{s: x.s, e: 0, d: []int32{0}} // trunc(0.5) = 0
	}
	return c.finalise(new(Decimal).Set(x), int(x.e)+1, RoundDown, false)
}

// Trunc rounds Decimal x towards zero using default context.
func (x *Decimal) Trunc() *Decimal {
	return globalContext.Trunc(x)
}

// Floor rounds Decimal x towards -Infinity to an integer.
func (c *Context) Floor(x *Decimal) *Decimal {
	if x == nil || !x.IsFinite() || x.IsZero() {
		return x
	}
	if x.e < 0 {
		if x.s < 0 {
			return &Decimal{s: -1, e: 0, d: []int32{1}} // floor(-0.5) = -1
		}
		return &Decimal{s: 1, e: 0, d: []int32{0}} // floor(0.5) = 0
	}
	return c.finalise(new(Decimal).Set(x), int(x.e)+1, RoundFloor, false)
}

// Floor rounds Decimal x using default context.
func (x *Decimal) Floor() *Decimal {
	return globalContext.Floor(x)
}

// Ceil rounds Decimal x towards +Infinity to an integer.
func (c *Context) Ceil(x *Decimal) *Decimal {
	if x == nil || !x.IsFinite() || x.IsZero() {
		return x
	}
	if x.e < 0 {
		if x.s < 0 {
			return &Decimal{s: -1, e: 0, d: []int32{0}} // ceil(-0.5) = -0
		}
		return &Decimal{s: 1, e: 0, d: []int32{1}} // ceil(0.5) = 1
	}
	return c.finalise(new(Decimal).Set(x), int(x.e)+1, RoundCeil, false)
}

// Ceil rounds Decimal x using default context.
func (x *Decimal) Ceil() *Decimal {
	return globalContext.Ceil(x)
}

// Round rounds Decimal x to an integer using configured rounding mode.
func (c *Context) Round(x *Decimal) *Decimal {
	if x == nil || !x.IsFinite() {
		return x
	}
	return c.finalise(new(Decimal).Set(x), int(x.e)+1, c.Rounding, false)
}

// Round rounds Decimal x using default context.
func (x *Decimal) Round() *Decimal {
	return globalContext.Round(x)
}
