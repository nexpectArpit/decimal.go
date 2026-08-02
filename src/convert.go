package decimal

import (
	"math"
	"strconv"
)

// Float64 converts Decimal x to a Go float64.
// Matches decimal.js toNumber() (lines 2205-2207).
func (x *Decimal) Float64() (float64, error) {
	if x == nil || x.s == 0 {
		return math.NaN(), nil
	}
	if x.d == nil {
		if x.s < 0 {
			return math.Inf(-1), nil
		}
		return math.Inf(1), nil
	}
	if x.IsZero() {
		if x.s < 0 {
			return math.Copysign(0.0, -1.0), nil
		}
		return 0.0, nil
	}

	str := x.String()
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToNumber is an alias for Float64.
func (x *Decimal) ToNumber() (float64, error) {
	return x.Float64()
}
