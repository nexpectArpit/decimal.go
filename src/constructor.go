package decimal

import (
	"fmt"
	"math"
	"strconv"
)

// New creates a new Decimal from string, int, int64, float64, or Decimal instance using the global context.
func New(v interface{}) (*Decimal, error) {
	return globalContext.New(v)
}

// New creates a new Decimal from string, int, int64, float64, or Decimal instance using the context c.
func (c *Context) New(v interface{}) (*Decimal, error) {
	if v == nil {
		return nil, fmt.Errorf("%w: nil", ErrInvalidArgument)
	}

	x := &Decimal{s: 1}

	switch val := v.(type) {
	case *Decimal:
		if val == nil {
			return nil, fmt.Errorf("%w: nil Decimal", ErrInvalidArgument)
		}
		x.s = val.s
		x.e = val.e
		if val.d != nil {
			x.d = make([]int32, len(val.d))
			copy(x.d, val.d)
		}
		return x, nil

	case int:
		return c.New(int64(val))

	case int64:
		if val < 0 {
			x.s = -1
			val = -val
		}
		return c.parseDecimal(x, strconv.FormatInt(val, 10)), nil

	case float64:
		if math.IsNaN(val) {
			x.s = 0
			x.d = nil
			x.e = 0
			return x, nil
		}
		if math.IsInf(val, 0) {
			if val < 0 {
				x.s = -1
			} else {
				x.s = 1
			}
			x.d = nil
			x.e = 0
			return x, nil
		}
		// Handle negative zero float64 (-0.0)
		if val == 0.0 {
			if 1/val < 0 || math.Signbit(val) {
				x.s = -1
			} else {
				x.s = 1
			}
			x.d = []int32{0}
			x.e = 0
			return x, nil
		}
		if val < 0 {
			x.s = -1
			val = -val
		}
		return c.parseDecimal(x, strconv.FormatFloat(val, 'g', -1, 64)), nil

	case string:
		str := val
		if len(str) == 0 {
			return nil, fmt.Errorf("%w: empty string", ErrInvalidArgument)
		}

		if str[0] == '-' {
			x.s = -1
			str = str[1:]
		} else if str[0] == '+' {
			str = str[1:]
		}

		if len(str) == 0 {
			return nil, fmt.Errorf("%w: invalid decimal string format", ErrInvalidArgument)
		}

		if isDecimal.MatchString(str) {
			return c.parseDecimal(x, str), nil
		}

		return c.parseOther(x, str)

	default:
		return nil, fmt.Errorf("%w: unsupported type %T", ErrInvalidArgument, v)
	}
}

// IsDecimal returns true if obj is a valid non-nil Decimal pointer.
func IsDecimal(obj interface{}) bool {
	d, ok := obj.(*Decimal)
	return ok && d != nil
}
