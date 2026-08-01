package decimal

import (
	"fmt"
	"strconv"
)

// finiteToString formats a finite Decimal into normal or exponential string representation.
// Matches decimal.js finiteToString() (lines 3113-3143).
func (c *Context) finiteToString(x *Decimal, isExp bool, sd int) string {
	if x == nil || x.s == 0 {
		return "NaN"
	}
	if x.d == nil {
		if x.s < 0 {
			return "-Infinity"
		}
		return "Infinity"
	}

	e := x.e
	str := digitsToString(x.d)
	length := len(str)
	var res string

	if isExp {
		if sd > 0 && sd > length {
			k := sd - length
			if length > 1 {
				res = str[:1] + "." + str[1:] + getZeroString(k)
			} else {
				res = str + "." + getZeroString(k)
			}
		} else if length > 1 {
			res = str[:1] + "." + str[1:]
		} else {
			res = str
		}

		expStr := strconv.FormatInt(e, 10)
		if e >= 0 {
			expStr = "+" + expStr
		}
		res = res + "e" + expStr
	} else if e < 0 {
		res = "0." + getZeroString(int(-e-1)) + str
		if sd > 0 && sd > length {
			res += getZeroString(sd - length)
		}
	} else if e >= int64(length) {
		res = str + getZeroString(int(e+1-int64(length)))
		if sd > 0 && int64(sd) > e+1 {
			res = res + "." + getZeroString(int(int64(sd)-e-1))
		}
	} else {
		k := int(e + 1)
		if k < length {
			res = str[:k] + "." + str[k:]
		} else {
			res = str
		}
		if sd > 0 && sd > length {
			if k == length {
				res += "."
			}
			res += getZeroString(sd - length)
		}
	}

	return res
}

// String formats the Decimal as a string using the global context rules.
// Matches decimal.js toString() (lines 2439-2445).
func (x *Decimal) String() string {
	return globalContext.String(x)
}

// String formats the Decimal as a string using context c rules.
func (c *Context) String(x *Decimal) string {
	if x == nil || x.s == 0 {
		return "NaN"
	}
	if x.d == nil {
		if x.s < 0 {
			return "-Infinity"
		}
		return "Infinity"
	}

	isExp := x.e <= c.ToExpNeg || x.e >= c.ToExpPos
	str := c.finiteToString(x, isExp, 0)

	if x.s < 0 && (len(x.d) == 0 || x.d[0] != 0) {
		return "-" + str
	}
	return str
}

// ToFixed returns a string representation in fixed-point notation with dp decimal places.
// Matches decimal.js toFixed() (lines 2026-2046).
func (c *Context) ToFixed(x *Decimal, dp int) string {
	if x == nil || !x.IsFinite() {
		return c.String(x)
	}

	y := c.finalise(new(Decimal).Set(x), int(x.e)+dp+1, c.Rounding, false)
	str := c.finiteToString(y, false, int(y.e)+dp+1)

	if x.IsNeg() && !x.IsZero() {
		return "-" + str
	}
	return str
}

// ToFixed returns a string representation in fixed-point notation using default context.
func (x *Decimal) ToFixed(dp int) string {
	return globalContext.ToFixed(x, dp)
}

// ToExponential returns a string representation in exponential notation with dp decimal places.
// Matches decimal.js toExponential() (lines 1730-1746).
func (c *Context) ToExponential(x *Decimal, dp int) string {
	if x == nil || !x.IsFinite() {
		return c.String(x)
	}

	y := c.finalise(new(Decimal).Set(x), dp+1, c.Rounding, false)
	str := c.finiteToString(y, true, dp+1)

	if x.IsNeg() && !x.IsZero() {
		return "-" + str
	}
	return str
}

// ToExponential returns a string representation in exponential notation using default context.
func (x *Decimal) ToExponential(dp int) string {
	return globalContext.ToExponential(x, dp)
}

// ToPrecision returns a string representation rounded to sd significant digits.
// Matches decimal.js toPrecision() (lines 2379-2397).
func (c *Context) ToPrecision(x *Decimal, sd int) string {
	if x == nil || !x.IsFinite() {
		return c.String(x)
	}

	y := c.finalise(new(Decimal).Set(x), sd, c.Rounding, false)
	isExp := sd <= int(y.e) || y.e <= c.ToExpNeg
	str := c.finiteToString(y, isExp, sd)

	if x.IsNeg() && !x.IsZero() {
		return "-" + str
	}
	return str
}

// ToPrecision returns a string representation rounded to sd significant digits using default context.
func (x *Decimal) ToPrecision(sd int) string {
	return globalContext.ToPrecision(x, sd)
}

// Set copies the content of src into x.
func (x *Decimal) Set(src *Decimal) *Decimal {
	if src == nil {
		x.s = 0
		x.e = 0
		x.d = nil
		return x
	}
	x.s = src.s
	x.e = src.e
	if src.d != nil {
		x.d = make([]int32, len(src.d))
		copy(x.d, src.d)
	} else {
		x.d = nil
	}
	return x
}

// ValueOf / MarshalJSON returns the string representation including negative zero sign.
func (x *Decimal) ValueOf() string {
	if x == nil || x.s == 0 {
		return "NaN"
	}
	if x.d == nil {
		if x.s < 0 {
			return "-Infinity"
		}
		return "Infinity"
	}
	isExp := x.e <= globalContext.ToExpNeg || x.e >= globalContext.ToExpPos
	str := globalContext.finiteToString(x, isExp, 0)
	if x.s < 0 {
		return "-" + str
	}
	return str
}

// MarshalJSON implements json.Marshaler.
func (x *Decimal) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", x.ValueOf())), nil
}
