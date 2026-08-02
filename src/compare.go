package decimal

// compare compares digit array a of length aL against b of length bL in base 10^7.
// Matches decimal.js compare() helper (lines 2697-2712).
func compareLimbs(a, b []int32, aL, bL int) int {
	minL := aL
	if bL < minL {
		minL = bL
	}

	for i := 0; i < minL; i++ {
		if a[i] != b[i] {
			if a[i] > b[i] {
				return 1
			}
			return -1
		}
	}

	if aL == bL {
		return 0
	}
	if aL > bL {
		return 1
	}
	return -1
}

func compareLimbsDiv(a, b []int32, aL, bL int) int {
	if aL != bL {
		if aL > bL {
			return 1
		}
		return -1
	}

	for i := 0; i < aL; i++ {
		if a[i] != b[i] {
			if a[i] > b[i] {
				return 1
			}
			return -1
		}
	}

	return 0
}

// Cmp compares Decimal x to Decimal y.
// Returns 1 if x > y, -1 if x < y, 0 if x == y, or 0 with false if either is NaN.
// Matches decimal.js comparedTo / cmp (lines 246-278).
func (x *Decimal) Cmp(y *Decimal) int {
	if x == nil || y == nil {
		return 0
	}

	xd := x.d
	yd := y.d
	xs := x.s
	ys := y.s

	// Either NaN or ±Infinity?
	if xd == nil || yd == nil {
		if xs == 0 || ys == 0 {
			return 0 // NaN
		}
		if xs != ys {
			return int(xs)
		}
		if xd == nil && yd == nil {
			return 0
		}
		if xd == nil {
			return int(xs)
		}
		return int(-ys)
	}

	// Either zero?
	if len(xd) > 0 && xd[0] == 0 || len(yd) > 0 && yd[0] == 0 {
		if len(xd) > 0 && xd[0] != 0 {
			return int(xs)
		}
		if len(yd) > 0 && yd[0] != 0 {
			return int(-ys)
		}
		return 0
	}

	// Signs differ?
	if xs != ys {
		return int(xs)
	}

	// Compare exponents
	if x.e != y.e {
		if (x.e > y.e) != (xs < 0) {
			return 1
		}
		return -1
	}

	xdL := len(xd)
	ydL := len(yd)
	minL := xdL
	if ydL < minL {
		minL = ydL
	}

	// Compare digit by digit
	for i := 0; i < minL; i++ {
		if xd[i] != yd[i] {
			if (xd[i] > yd[i]) != (xs < 0) {
				return 1
			}
			return -1
		}
	}

	// Compare lengths
	if xdL == ydL {
		return 0
	}
	if (xdL > ydL) != (xs < 0) {
		return 1
	}
	return -1
}

// Cmp compares Decimal x to Decimal y using context c settings.
func (c *Context) Cmp(x, y *Decimal) int {
	return x.Cmp(y)
}

// Eq returns true if x equals y.
func (x *Decimal) Eq(y *Decimal) bool {
	if x.IsNaN() || y.IsNaN() {
		return false
	}
	return x.Cmp(y) == 0
}

// Gt returns true if x is greater than y.
func (x *Decimal) Gt(y *Decimal) bool {
	if x.IsNaN() || y.IsNaN() {
		return false
	}
	return x.Cmp(y) > 0
}

// Gte returns true if x is greater than or equal to y.
func (x *Decimal) Gte(y *Decimal) bool {
	if x.IsNaN() || y.IsNaN() {
		return false
	}
	return x.Cmp(y) >= 0
}

// Lt returns true if x is less than y.
func (x *Decimal) Lt(y *Decimal) bool {
	if x.IsNaN() || y.IsNaN() {
		return false
	}
	return x.Cmp(y) < 0
}

// Lte returns true if x is less than or equal to y.
func (x *Decimal) Lte(y *Decimal) bool {
	if x.IsNaN() || y.IsNaN() {
		return false
	}
	res := x.Cmp(y)
	return res == -1 || res == 0
}

// IsNaN returns true if Decimal x is NaN.
func (x *Decimal) IsNaN() bool {
	return x == nil || x.s == 0
}

// IsFinite returns true if Decimal x is finite (neither NaN nor ±Infinity).
func (x *Decimal) IsFinite() bool {
	return x != nil && x.d != nil && x.s != 0
}

// IsZero returns true if Decimal x represents zero.
func (x *Decimal) IsZero() bool {
	return x != nil && len(x.d) > 0 && x.d[0] == 0
}

// IsNeg returns true if Decimal x is negative.
func (x *Decimal) IsNeg() bool {
	return x != nil && x.s < 0
}

// IsPos returns true if Decimal x is positive.
func (x *Decimal) IsPos() bool {
	return x != nil && x.s > 0
}

// IsInt returns true if Decimal x is an integer.
func (x *Decimal) IsInt() bool {
	if x == nil || !x.IsFinite() {
		return false
	}
	if len(x.d) == 0 || x.d[0] == 0 {
		return true
	}
	if x.e < 0 {
		return false
	}
	expLimbs := int(x.e / int64(LogBase))
	if expLimbs >= len(x.d)-1 {
		return true
	}
	for i := expLimbs + 1; i < len(x.d); i++ {
		if x.d[i] != 0 {
			return false
		}
	}
	return true
}
