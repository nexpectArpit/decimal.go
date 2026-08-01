package decimal

import (
	"errors"
)

// Constants matching decimal.js specification
const (
	// Base is the radix for the coefficient array d (10^7).
	Base int32 = 10000000

	// LogBase is the number of decimal digits per limb element.
	LogBase int = 7

	// ExpLimit is the maximum exponent magnitude (9e15).
	ExpLimit int64 = 9000000000000000

	// MaxDigits is the upper limit on precision and formatting digit arguments (1e9).
	MaxDigits int = 1000000000

	// MaxSafeInteger is the upper bound for exact integer exponentiation (2^53 - 1).
	MaxSafeInteger int64 = 9007199254740991

	// Numerals is the alphabet used for base conversion.
	Numerals string = "0123456789abcdef"
)

// RoundingMode defines the 9 IEEE-754 and custom rounding modes supported by decimal.js.
type RoundingMode uint8

const (
	// RoundUp rounds away from zero (0).
	RoundUp RoundingMode = 0

	// RoundDown rounds towards zero / truncation (1).
	RoundDown RoundingMode = 1

	// RoundCeil rounds towards +Infinity (2).
	RoundCeil RoundingMode = 2

	// RoundFloor rounds towards -Infinity (3).
	RoundFloor RoundingMode = 3

	// RoundHalfUp rounds towards nearest neighbor; if equidistant, rounds up (4).
	RoundHalfUp RoundingMode = 4

	// RoundHalfDown rounds towards nearest neighbor; if equidistant, rounds down (5).
	RoundHalfDown RoundingMode = 5

	// RoundHalfEven rounds towards nearest neighbor; if equidistant, rounds towards even neighbor (6).
	RoundHalfEven RoundingMode = 6

	// RoundHalfCeil rounds towards nearest neighbor; if equidistant, rounds towards +Infinity (7).
	RoundHalfCeil RoundingMode = 7

	// RoundHalfFloor rounds towards nearest neighbor; if equidistant, rounds towards -Infinity (8).
	RoundHalfFloor RoundingMode = 8
)

// ModuloMode defines the modulo calculation mode for a mod n operations.
type ModuloMode uint8

const (
	// ModuloUp: Remainder is positive if dividend is negative, else negative (0).
	ModuloUp ModuloMode = 0

	// ModuloTrunc: Remainder has same sign as dividend, JS % (1).
	ModuloTrunc ModuloMode = 1

	// ModuloFloor: Remainder has same sign as divisor, Python % (3).
	ModuloFloor ModuloMode = 3

	// ModuloHalfEven: IEEE 754 remainder function (6).
	ModuloHalfEven ModuloMode = 6

	// ModuloEuclid: Euclidean division, remainder is always positive (9).
	ModuloEuclid ModuloMode = 9
)

// Sentinel error definitions matching decimal.js Error messages
var (
	ErrInvalidArgument        = errors.New("[DecimalError] Invalid argument")
	ErrPrecisionLimitExceeded = errors.New("[DecimalError] Precision limit exceeded")
	ErrCryptoUnavailable      = errors.New("[DecimalError] crypto unavailable")
	ErrDivisionByZero         = errors.New("[DecimalError] Division by zero")
)

// Decimal represents an arbitrary-precision decimal number.
// Its observable behavior matches decimal.js exactly.
type Decimal struct {
	s int8    // Sign: 1 for positive, -1 for negative, 0 for NaN/Zero sentinel
	e int64   // Exponent: base-10 exponent of the number
	d []int32 // Coefficients: limb array in radix 10^7. Nil if NaN or ±Infinity.
}

// Sign returns the sign of the Decimal (1, -1, or 0).
func (d *Decimal) Sign() int8 {
	return d.s
}

// Exponent returns the base-10 exponent of the Decimal.
func (d *Decimal) Exponent() int64 {
	return d.e
}

// Coefficients returns a copy of the coefficient slice.
func (d *Decimal) Coefficients() []int32 {
	if d.d == nil {
		return nil
	}
	res := make([]int32, len(d.d))
	copy(res, d.d)
	return res
}
