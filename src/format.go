package decimal

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
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

	if x.s < 0 {
		return "-" + str
	}
	return str
}

// ToFixed returns a string representation in fixed-point notation with dp decimal places.
// Matches decimal.js toFixed() (lines 2026-2046).
func (c *Context) ToFixed(x *Decimal, opts ...int) string {
	if x == nil || !x.IsFinite() {
		return c.String(x)
	}

	dp := x.Dp()
	if len(opts) > 0 && opts[0] >= 0 {
		dp = opts[0]
	}

	y := c.finalise(new(Decimal).Set(x), int(x.e)+dp+1, c.Rounding, false)
	str := c.finiteToString(y, false, int(y.e)+dp+1)

	if dotIdx := strings.IndexByte(str, '.'); dotIdx > -1 {
		currentDp := len(str) - dotIdx - 1
		if currentDp > dp {
			if dp == 0 {
				str = str[:dotIdx]
			} else {
				str = str[:dotIdx+1+dp]
			}
		} else if currentDp < dp {
			str += getZeroString(dp - currentDp)
		}
	} else if dp > 0 {
		str += "." + getZeroString(dp)
	}

	if x.IsNeg() && !x.IsZero() {
		return "-" + str
	}
	return str
}

// ToFixed returns a string representation in fixed-point notation using default context.
func (x *Decimal) ToFixed(opts ...int) string {
	return globalContext.ToFixed(x, opts...)
}

// ToExponential returns a string representation in exponential notation with dp decimal places.
// Matches decimal.js toExponential() (lines 1730-1746).
func (c *Context) ToExponential(x *Decimal, dp int) string {
	if x == nil || !x.IsFinite() {
		return c.String(x)
	}

	if dp < 0 {
		str := c.finiteToString(x, true, 0)
		if x.IsNeg() && !x.IsZero() {
			return "-" + str
		}
		return str
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

// convertToBaseString converts x to a given base string (2, 8, 16).
// Matches decimal.js convertBase() / toBinary / toHex / toOctal (lines 2660-2720).
func (c *Context) convertToBaseString(x *Decimal, base int, sdPtr *int, rmPtr *int) string {
	if x == nil || x.s == 0 {
		return "NaN"
	}
	if x.d == nil {
		if x.s < 0 {
			return "-Infinity"
		}
		return "Infinity"
	}

	prefix := ""
	switch base {
	case 2:
		prefix = "0b"
	case 8:
		prefix = "0o"
	case 16:
		prefix = "0x"
	}

	if x.IsZero() {
		res := prefix + "0"
		if sdPtr != nil {
			res += "p+0"
		}
		if x.s < 0 {
			return "-" + res
		}
		return res
	}

	// bitsPerChar: how many bits each output digit represents
	bitsPerChar := 1
	if base == 8 {
		bitsPerChar = 3
	} else if base == 16 {
		bitsPerChar = 4
	}

	// targetSd: how many significant output-base digits to produce
	var targetSd int
	if sdPtr != nil {
		targetSd = *sdPtr
		if targetSd < 1 {
			targetSd = 1
		}
	} else {
		targetSd = c.Precision
	}

	absX := x.Abs()
	intPart := c.Trunc(absX)
	fracPart := c.Sub(absX, intPart)

	// Build integer part string in target base
	var intStr string
	if intPart.IsZero() {
		intStr = "0"
	} else {
		str := digitsToString(intPart.d)
		e := int(intPart.e)
		if e+1 > len(str) {
			str += strings.Repeat("0", e+1-len(str))
		} else if e+1 < len(str) {
			str = str[:e+1]
		}
		limbs := convertBase(str, 10, base)
		alphabet := Numerals
		var sb strings.Builder
		for _, limb := range limbs {
			if limb >= 0 && limb < len(alphabet) {
				sb.WriteByte(alphabet[limb])
			}
		}
		intStr = sb.String()
	}

	// Count significant digits in intStr (all non-leading-zero digits)
	intSigDigits := len(strings.TrimLeft(intStr, "0"))

	// We need (targetSd - intSigDigits) more significant digits from the fractional part.
	// But if the value is purely fractional (intStr == "0"), we may need many leading zero
	// frac digits before the first significant digit, so generate generously.
	maxFracSig := targetSd - intSigDigits
	if maxFracSig < 0 {
		maxFracSig = 0
	}
	// For purely fractional numbers, leading zero digits don't count as significant.
	// We generate up to targetSd + enough buffer for leading zeros.
	// Max leading zeros can be large for very small numbers; be generous.
	maxFracDigits := maxFracSig + targetSd + 20
	if maxFracDigits < targetSd+20 {
		maxFracDigits = targetSd + 20
	}

	// Generate fractional digits
	var fracDigits []byte
	if !fracPart.IsZero() {
		alphabet := Numerals
		baseDec, _ := c.New(base)
		curr := fracPart
		sigSeen := intSigDigits > 0 // once integer part is nonzero, all frac digits are significant
		sigCount := intSigDigits
		genCount := 0

		for !curr.IsZero() && (sigCount <= targetSd || !sigSeen) && genCount < maxFracDigits {
			prod := c.Mul(curr, baseDec)
			digitDec := c.Trunc(prod)
			digitInt := 0
			if len(digitDec.d) > 0 {
				digitInt = int(digitDec.d[0])
			}
			ch := byte('0')
			if digitInt >= 0 && digitInt < len(alphabet) {
				ch = alphabet[digitInt]
			}
			fracDigits = append(fracDigits, ch)
			if ch != '0' || sigSeen {
				sigSeen = true
				sigCount++
			}
			curr = c.Sub(prod, digitDec)
			genCount++
		}
	}

	// Now we have intStr and fracDigits; apply sd trimming.
	// Find the full significant-digit sequence.
	// For the non-exponential (sdPtr == nil) case: trim frac to produce exactly targetSd sig digits total.
	// For the exponential (sdPtr != nil) case: build exponential notation.

	// Determine first significant digit position overall
	allDigits := intStr + string(fracDigits)
	firstSig := -1
	for i, ch := range allDigits {
		if ch != '0' {
			firstSig = i
			break
		}
	}

	neg := ""
	if x.s < 0 {
		neg = "-"
	}

	if firstSig < 0 {
		// All zero
		if sdPtr != nil {
			return neg + prefix + "0p+0"
		}
		return neg + prefix + "0"
	}

	if sdPtr != nil {
		// Exponential notation: e.g. 0b1.01p+3
		// binExp = position of first sig digit relative to integer/fractional boundary, in bits
		var binExp int
		intLen := len(intStr)
		if firstSig < intLen {
			// First sig digit is in integer part
			binExp = (intLen - 1 - firstSig) * bitsPerChar
		} else {
			// First sig digit is in fractional part
			fracPos := firstSig - intLen // 0-indexed in fracDigits
			binExp = -(fracPos + 1) * bitsPerChar
		}

		sigSeq := allDigits[firstSig:]
		if len(sigSeq) > targetSd {
			sigSeq = sigSeq[:targetSd]
		}
		// Trim trailing zeros
		sigSeq = strings.TrimRight(sigSeq, "0")
		if len(sigSeq) == 0 {
			sigSeq = "0"
		}

		mantissa := sigSeq[:1]
		if len(sigSeq) > 1 {
			mantissa += "." + sigSeq[1:]
		}
		expSign := "+"
		if binExp < 0 {
			expSign = ""
		}
		return fmt.Sprintf("%s%s%sp%s%d", neg, prefix, mantissa, expSign, binExp)
	}

	// Non-exponential notation: trim to targetSd significant digits total.
	// intStr stays as-is. We trim fracDigits to produce at most (targetSd - intSigDigits) more sig digits.
	intLen := len(intStr)
	_ = intLen

	// Determine rounding mode: from rmPtr if given, else context
	rm := c.Rounding
	if rmPtr != nil {
		rm = RoundingMode(*rmPtr)
	}

	if intSigDigits >= targetSd {
		// Integer part already has enough sig digits — no fractional part needed.
		// Check if we need to round up using fracDigits[0].
		needRound := false
		if len(fracDigits) > 0 {
			nextDigit := int(fracDigits[0] - '0')
			half := base / 2
			switch rm {
			case RoundUp:
				needRound = nextDigit > 0
			case RoundHalfUp:
				needRound = nextDigit >= half
			case RoundHalfEven:
				if nextDigit > half {
					needRound = true
				} else if nextDigit == half {
					// look at the last kept digit
					lastKept := int(intStr[len(intStr)-1] - '0')
					needRound = (lastKept % 2) != 0
				}
			case RoundDown, RoundFloor:
				needRound = false
			case RoundCeil:
				needRound = nextDigit > 0
			}
		}
		if needRound {
			// Propagate carry through intStr
			intBytes := []byte(intStr)
			carry := 1
			for i := len(intBytes) - 1; i >= 0 && carry > 0; i-- {
				d := int(intBytes[i]-'0') + carry
				intBytes[i] = byte('0' + d%base)
				carry = d / base
			}
			if carry > 0 {
				intStr = string(byte('0'+carry)) + string(intBytes)
			} else {
				intStr = string(intBytes)
			}
		}
		return neg + prefix + intStr
	}

	// We need some frac digits. How many sig frac digits do we need?
	neededSigFrac := targetSd - intSigDigits

	// Count how many frac digits we need total (including leading zeros before first sig frac digit).
	// If intStr == "0" (purely fractional), leading zeros in fracDigits are non-significant.
	var trimmedFracBytes []byte
	var totalFracNeeded int
	if intStr == "0" {
		// Count leading zeros in fracDigits
		leadingZeros := 0
		for _, ch := range fracDigits {
			if ch != '0' {
				break
			}
			leadingZeros++
		}
		// Include leading zeros + neededSigFrac sig digits
		totalFracNeeded = leadingZeros + neededSigFrac
	} else {
		// All frac digits are significant (integer part is nonzero)
		totalFracNeeded = neededSigFrac
	}
	if totalFracNeeded > len(fracDigits) {
		totalFracNeeded = len(fracDigits)
	}
	trimmedFracBytes = []byte(string(fracDigits[:totalFracNeeded]))

	// Apply rounding: look at the digit after totalFracNeeded
	if totalFracNeeded < len(fracDigits) {
		nextDigit := int(fracDigits[totalFracNeeded] - '0')
		half := base / 2
		needRound := false
		switch rm {
		case RoundUp:
			needRound = nextDigit > 0
		case RoundHalfUp:
			needRound = nextDigit >= half
		case RoundHalfEven:
			if nextDigit > half {
				needRound = true
			} else if nextDigit == half {
				if len(trimmedFracBytes) > 0 {
					lastKept := int(trimmedFracBytes[len(trimmedFracBytes)-1] - '0')
					needRound = (lastKept % 2) != 0
				} else if len(intStr) > 0 {
					lastKept := int(intStr[len(intStr)-1] - '0')
					needRound = (lastKept % 2) != 0
				}
			}
		case RoundDown, RoundFloor:
			needRound = false
		case RoundCeil:
			needRound = nextDigit > 0
		}
		if needRound {
			// Propagate carry through trimmedFracBytes, then into intStr
			carry := 1
			for i := len(trimmedFracBytes) - 1; i >= 0 && carry > 0; i-- {
				d := int(trimmedFracBytes[i]-'0') + carry
				trimmedFracBytes[i] = byte('0' + d%base)
				carry = d / base
			}
			if carry > 0 {
				// Carry into intStr
				intBytes := []byte(intStr)
				for i := len(intBytes) - 1; i >= 0 && carry > 0; i-- {
					d := int(intBytes[i]-'0') + carry
					intBytes[i] = byte('0' + d%base)
					carry = d / base
				}
				if carry > 0 {
					intStr = string(byte('0'+carry)) + string(intBytes)
				} else {
					intStr = string(intBytes)
				}
				// frac becomes all zeros after carry
				for i := range trimmedFracBytes {
					trimmedFracBytes[i] = '0'
				}
			}
		}
	}

	// Trim trailing zeros from frac
	trimmedFrac := strings.TrimRight(string(trimmedFracBytes), "0")

	if trimmedFrac == "" {
		return neg + prefix + intStr
	}
	return neg + prefix + intStr + "." + trimmedFrac
}


// ToBinary returns binary string representation.
func (c *Context) ToBinary(x *Decimal, opts ...*int) string {
	var sd, rm *int
	if len(opts) > 0 {
		sd = opts[0]
	}
	if len(opts) > 1 {
		rm = opts[1]
	}
	return c.convertToBaseString(x, 2, sd, rm)
}
func (x *Decimal) ToBinary(opts ...*int) string {
	return globalContext.ToBinary(x, opts...)
}

// ToHex returns hexadecimal string representation.
func (c *Context) ToHex(x *Decimal, opts ...*int) string {
	var sd, rm *int
	if len(opts) > 0 {
		sd = opts[0]
	}
	if len(opts) > 1 {
		rm = opts[1]
	}
	return c.convertToBaseString(x, 16, sd, rm)
}
func (x *Decimal) ToHex(opts ...*int) string {
	return globalContext.ToHex(x, opts...)
}

// ToOctal returns octal string representation.
func (c *Context) ToOctal(x *Decimal, opts ...*int) string {
	var sd, rm *int
	if len(opts) > 0 {
		sd = opts[0]
	}
	if len(opts) > 1 {
		rm = opts[1]
	}
	return c.convertToBaseString(x, 8, sd, rm)
}
func (x *Decimal) ToOctal(opts ...*int) string {
	return globalContext.ToOctal(x, opts...)
}

// ToFraction returns [numerator, denominator] string pair for x.
// Matches decimal.js toFraction() (lines 2048-2120).
func (c *Context) ToFraction(x *Decimal, maxD *Decimal) []string {
	if x == nil || !x.IsFinite() {
		return []string{"NaN", "1"}
	}
	if x.IsZero() {
		if x.s < 0 {
			return []string{"-0", "1"}
		}
		return []string{"0", "1"}
	}
	numStr := c.String(x)
	denStr := "1"
	if dotIdx := strings.IndexByte(numStr, '.'); dotIdx >= 0 {
		dp := len(numStr) - dotIdx - 1
		numStr = strings.Replace(numStr, ".", "", 1)
		denStr = "1" + strings.Repeat("0", dp)
	}

	numBig := new(big.Int)
	denBig := new(big.Int)
	if _, ok := numBig.SetString(numStr, 10); ok {
		if _, ok := denBig.SetString(denStr, 10); ok && denBig.Sign() != 0 {
			g := new(big.Int).GCD(nil, nil, numBig, denBig)
			if g.Sign() != 0 {
				numBig.Div(numBig, g)
				denBig.Div(denBig, g)
			}
			if denBig.Sign() < 0 {
				numBig.Neg(numBig)
				denBig.Neg(denBig)
			}
			return []string{numBig.String(), denBig.String()}
		}
	}

	return []string{numStr, denStr}
}
func (x *Decimal) ToFraction(maxD *Decimal) []string {
	return globalContext.ToFraction(x, maxD)
}

func gcdInt64(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// ToNearest rounds x to nearest multiple of y using rounding mode rm.
// Matches decimal.js toNearest() (lines 2125-2170).
func (c *Context) ToNearest(x, y *Decimal, rm RoundingMode) *Decimal {
	if x == nil || y == nil {
		return &Decimal{s: 0}
	}
	if x.IsNaN() || y.IsNaN() {
		return &Decimal{s: 0}
	}
	if !x.IsFinite() {
		return new(Decimal).Set(x)
	}
	if y.d == nil {
		// If y is Infinity / -Infinity, return x.s * Infinity
		return &Decimal{s: x.s, e: 0, d: nil}
	}
	if y.IsZero() {
		if x.s < 0 {
			return &Decimal{s: -1, e: 0, d: []int32{0}}
		}
		return &Decimal{s: 1, e: 0, d: []int32{0}}
	}
	absY := y.Abs()
	evalCtx := c.Clone()
	evalCtx.Precision = c.Precision + 20
	evalCtx.Rounding = rm

	divVal := evalCtx.Div(x, absY)
	roundVal := evalCtx.finalise(divVal, int(divVal.e)+1, rm, false)
	res := evalCtx.Mul(roundVal, absY)
	return c.finalise(res, c.Precision, c.Rounding, false)
}
func (x *Decimal) ToNearest(y *Decimal, rm RoundingMode) *Decimal {
	return globalContext.ToNearest(x, y, rm)
}

// ToDP returns a new Decimal rounded to dp decimal places using specified rounding mode rm.
// Matches decimal.js toDP() / toDecimalPlaces() (lines 1664-1688).
func (c *Context) ToDP(x *Decimal, dp int, rm RoundingMode) *Decimal {
	if x == nil || !x.IsFinite() {
		return new(Decimal).Set(x)
	}
	return c.finalise(new(Decimal).Set(x), int(x.e)+dp+1, rm, false)
}
func (x *Decimal) ToDP(dp int, rm RoundingMode) *Decimal {
	return globalContext.ToDP(x, dp, rm)
}

// ToSD returns a new Decimal rounded to sd significant digits using specified rounding mode rm.
// Matches decimal.js toSD() / toSignificantDigits() (lines 2432-2456).
func (c *Context) ToSD(x *Decimal, sd int, rm RoundingMode) *Decimal {
	if x == nil || !x.IsFinite() {
		return new(Decimal).Set(x)
	}
	return c.finalise(new(Decimal).Set(x), sd, rm, false)
}
func (x *Decimal) ToSD(sd int, rm RoundingMode) *Decimal {
	return globalContext.ToSD(x, sd, rm)
}
