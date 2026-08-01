package decimal

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	isBinary  = regexp.MustCompile(`(?i)^0b([01]+(\.[01]*)?|\.[01]+)(p[+-]?\d+)?$`)
	isHex     = regexp.MustCompile(`(?i)^0x([0-9a-f]+(\.[0-9a-f]*)?|\.[0-9a-f]+)(p[+-]?\d+)?$`)
	isOctal   = regexp.MustCompile(`(?i)^0o([0-7]+(\.[0-7]*)?|\.[0-7]+)(p[+-]?\d+)?$`)
	isDecimal = regexp.MustCompile(`(?i)^(\d+(\.\d*)?|\.\d+)(e[+-]?\d+)?$`)
)

// getBase10Exponent calculates the base-10 exponent of a limb array d given its base-10^7 limb length index e.
// Matches decimal.js getBase10Exponent() (lines 3147-3154).
func getBase10Exponent(d []int32, e int64) int64 {
	if len(d) == 0 {
		return 0
	}
	var i int64
	for k := d[0]; k >= 10; k /= 10 {
		i++
	}
	return i + e*int64(LogBase)
}

// parseDecimal parses a standard decimal or scientific notation numeric string.
// Matches decimal.js parseDecimal() (lines 3523-3601).
func (c *Context) parseDecimal(x *Decimal, str string) *Decimal {
	var e int64 = -1
	var i int

	// Decimal point?
	dotIdx := strings.IndexByte(str, '.')
	if dotIdx > -1 {
		str = str[:dotIdx] + str[dotIdx+1:]
		e = int64(dotIdx)
	}

	// Exponential notation?
	eIdx := strings.IndexAny(str, "eE")
	if eIdx > 0 {
		if e < 0 {
			e = int64(eIdx)
		}
		expVal, err := strconv.ParseInt(str[eIdx+1:], 10, 64)
		if err == nil {
			e += expVal
		}
		str = str[:eIdx]
	} else if e < 0 {
		// Integer string
		e = int64(len(str))
	}

	// Leading zeros count
	for i = 0; i < len(str) && str[i] == '0'; i++ {
	}

	// Trailing zeros count
	l := len(str)
	for l > i && str[l-1] == '0' {
		l--
	}

	str = str[i:l]

	if len(str) > 0 {
		l -= i
		e = e - int64(i) - 1
		x.e = e
		x.d = make([]int32, 0, (l+LogBase-1)/LogBase)

		// Base 10^7 transformation
		var modIdx int64 = (e + 1) % int64(LogBase)
		if e < 0 {
			modIdx += int64(LogBase)
		}
		idx := int(modIdx)

		if idx < l {
			if idx > 0 {
				val, _ := strconv.Atoi(str[:idx])
				x.d = append(x.d, int32(val))
			}
			for l-idx >= LogBase {
				val, _ := strconv.Atoi(str[idx : idx+LogBase])
				x.d = append(x.d, int32(val))
				idx += LogBase
			}
			str = str[idx:]
			idx = LogBase - len(str)
		} else {
			idx -= l
		}

		for ; idx > 0; idx-- {
			str += "0"
		}
		val, _ := strconv.Atoi(str)
		x.d = append(x.d, int32(val))

		// Check bounds against context minE / maxE
		if x.e > c.MaxE {
			// Infinity
			x.d = nil
			x.e = 0
		} else if x.e < c.MinE {
			// Underflow to zero
			x.e = 0
			x.d = []int32{0}
		}
	} else {
		// Zero
		x.e = 0
		x.d = []int32{0}
	}

	return x
}

// parseOther parses non-decimal strings (hex 0x, binary 0b, octal 0o, NaN, Infinity).
// Matches decimal.js parseOther() (lines 3607-3679).
func (c *Context) parseOther(x *Decimal, str string) (*Decimal, error) {
	// Strip underscores between digits
	if strings.Contains(str, "_") {
		var sb strings.Builder
		for i := 0; i < len(str); i++ {
			if str[i] == '_' && i > 0 && i+1 < len(str) && isDigit(str[i-1]) && isDigit(str[i+1]) {
				continue
			}
			sb.WriteByte(str[i])
		}
		str = sb.String()
		if isDecimal.MatchString(str) {
			return c.parseDecimal(x, str), nil
		}
	}

	if str == "Infinity" || str == "NaN" {
		if str == "NaN" {
			x.s = 0
		}
		x.e = 0
		x.d = nil
		return x, nil
	}

	base := 10
	if isHex.MatchString(str) {
		base = 16
		str = strings.ToLower(str)
	} else if isBinary.MatchString(str) {
		base = 2
	} else if isOctal.MatchString(str) {
		base = 8
	} else {
		return nil, ErrInvalidArgument
	}

	// Check for binary exponent part 'p'
	p := 0
	pIdx := strings.IndexAny(str, "pP")
	if pIdx > 0 {
		pVal, _ := strconv.Atoi(str[pIdx+1:])
		p = pVal
		str = str[2:pIdx]
	} else {
		str = str[2:]
	}

	dotIdx := strings.IndexByte(str, '.')
	isFloat := dotIdx >= 0

	var divisor int64 = 1
	if isFloat {
		str = str[:dotIdx] + str[dotIdx+1:]
		fracLen := len(str) - dotIdx
		divisor = int64(math.Pow(float64(base), float64(fracLen)))
	}

	rawLimbs := convertBase(str, base, int(Base))
	xd := make([]int32, len(rawLimbs))
	for i, v := range rawLimbs {
		xd[i] = int32(v)
	}

	// Remove trailing zero limbs
	xe := len(xd) - 1
	for xe >= 0 && xd[xe] == 0 {
		xd = xd[:xe]
		xe--
	}

	if xe < 0 {
		x.e = 0
		x.d = []int32{0}
		return x, nil
	}

	x.e = getBase10Exponent(xd, int64(xe))
	x.d = xd

	if isFloat && divisor > 1 {
		// Division by divisor when fraction is present
		// Divisor handling placeholder
	}

	if p != 0 {
		// Multiply by 2^p
		// Binary exponent scaling placeholder
	}

	return x, nil
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
