package decimal

import (
	"fmt"
	"strconv"
	"strings"
)

// getZeroString returns a string of k zeros.
func getZeroString(k int) string {
	if k <= 0 {
		return ""
	}
	return strings.Repeat("0", k)
}

// checkInt32 validates that integer i is within inclusive range [min, max].
func checkInt32(i, min, max int) error {
	if i < min || i > max {
		return fmt.Errorf("%w: %d", ErrInvalidArgument, i)
	}
	return nil
}

// isOdd returns true if n is an odd integer.
func isOdd(n int64) bool {
	return n%2 != 0
}

// digitsToString converts a base 10^7 limb array d to its exact decimal digit string representation.
// Matches decimal.js digitsToString() (lines 2520-2547).
func digitsToString(d []int32) string {
	if len(d) == 0 {
		return "0"
	}

	indexOfLastWord := len(d) - 1
	w := d[0]

	if indexOfLastWord > 0 {
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(int(w)))

		for i := 1; i < indexOfLastWord; i++ {
			ws := strconv.Itoa(int(d[i]))
			k := LogBase - len(ws)
			if k > 0 {
				sb.WriteString(getZeroString(k))
			}
			sb.WriteString(ws)
		}

		w = d[indexOfLastWord]
		k := LogBase - len(strconv.Itoa(int(w)))
		if k > 0 {
			sb.WriteString(getZeroString(k))
		}

		// Remove trailing zeros of last word.
		for w%10 == 0 && w != 0 {
			w /= 10
		}
		sb.WriteString(strconv.Itoa(int(w)))
		return sb.String()
	}

	if w == 0 {
		return "0"
	}

	// Remove trailing zeros of single word
	for w%10 == 0 && w != 0 {
		w /= 10
	}
	return strconv.Itoa(int(w))
}

// convertBase converts a numeric string from baseIn to an array of integers in baseOut.
// Matches decimal.js convertBase() (lines 2613-2634).
func convertBase(str string, baseIn, baseOut int) []int {
	if len(str) == 0 {
		return []int{0}
	}

	arr := []int{0}

	for i := 0; i < len(str); i++ {
		for arrL := len(arr) - 1; arrL >= 0; arrL-- {
			arr[arrL] *= baseIn
		}

		charVal := strings.IndexByte(Numerals, str[i])
		if charVal < 0 {
			charVal = strings.IndexByte(strings.ToUpper(Numerals), str[i])
		}
		if charVal >= 0 {
			arr[0] += charVal
		}

		for j := 0; j < len(arr); j++ {
			if arr[j] >= baseOut {
				if j+1 == len(arr) {
					arr = append(arr, 0)
				}
				arr[j+1] += arr[j] / baseOut
				arr[j] %= baseOut
			}
		}
	}

	// Reverse array to maintain big-endian limb ordering
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}

	return arr
}

// Dp returns the number of decimal places (digits after decimal point).
// Matches decimal.js dp() (lines 1640-1662).
func (x *Decimal) Dp() int {
	if x == nil || !x.IsFinite() {
		return -1
	}
	if x.d == nil || len(x.d) == 0 {
		return 0
	}
	str := digitsToString(x.d)
	e := x.e
	length := int64(len(str))
	if e >= length-1 {
		return 0
	}
	if e < 0 {
		return int(-e - 1 + length)
	}
	return int(length - 1 - e)
}

// Sd returns the number of significant digits.
// Matches decimal.js sd() (lines 2404-2430).
func (x *Decimal) Sd(includeZeros ...bool) int {
	if x == nil || !x.IsFinite() {
		return -1
	}
	if x.d == nil || len(x.d) == 0 {
		return 0
	}
	str := digitsToString(x.d)
	length := len(str)
	if len(includeZeros) > 0 && includeZeros[0] && x.e+1 > int64(length) {
		return int(x.e + 1)
	}
	return length
}

// Clamp clamps x to range [min, max].
// Matches decimal.js clamp() (lines 1420-1440).
func (c *Context) Clamp(x, min, max *Decimal) *Decimal {
	if x == nil || min == nil || max == nil || x.IsNaN() || min.IsNaN() || max.IsNaN() {
		return &Decimal{s: 0}
	}
	if min.Gt(max) {
		return &Decimal{s: 0}
	}
	if x.Lt(min) {
		return c.finalise(new(Decimal).Set(min), c.Precision, c.Rounding, false)
	}
	if x.Gt(max) {
		return c.finalise(new(Decimal).Set(max), c.Precision, c.Rounding, false)
	}
	return c.finalise(new(Decimal).Set(x), c.Precision, c.Rounding, false)
}

func (x *Decimal) Clamp(min, max *Decimal) *Decimal {
	return globalContext.Clamp(x, min, max)
}
