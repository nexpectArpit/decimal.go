package decimal

import (
	"sync"
)

// Context holds precision, rounding, modulo, and exponent limit settings.
// It translates the JS closure state created by Decimal.clone() and Decimal.config().
type Context struct {
	mu        sync.RWMutex
	Precision int          // 1 to MaxDigits (default 20)
	Rounding  RoundingMode // 0 to 8 (default 4 = RoundHalfUp)
	Modulo    ModuloMode   // 0 to 9 (default 1 = ModuloTrunc)
	ToExpNeg  int64        // 0 to -ExpLimit (default -7)
	ToExpPos  int64        // 0 to ExpLimit (default 21)
	MinE      int64        // -1 to -ExpLimit (default -ExpLimit)
	MaxE      int64        // 1 to ExpLimit (default ExpLimit)
	Crypto    bool         // use secure crypto random (default false)
}

// DefaultContext creates a new Context populated with decimal.js default settings.
func DefaultContext() *Context {
	return &Context{
		Precision: 20,
		Rounding:  RoundHalfUp,
		Modulo:    ModuloTrunc,
		ToExpNeg:  -7,
		ToExpPos:  21,
		MinE:      -ExpLimit,
		MaxE:      ExpLimit,
		Crypto:    false,
	}
}

// Global Default Context matching DEFAULTS in decimal.js
var globalContext = DefaultContext()

// GlobalContext returns a copy of the active package-level default Context.
func GlobalContext() *Context {
	return globalContext.Clone()
}

// SetGlobalContext sets the active package-level default Context.
func SetGlobalContext(ctx *Context) {
	if ctx != nil {
		globalContext = ctx.Clone()
	}
}

// Clone creates an isolated copy of the Context.
func (c *Context) Clone() *Context {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return &Context{
		Precision: c.Precision,
		Rounding:  c.Rounding,
		Modulo:    c.Modulo,
		ToExpNeg:  c.ToExpNeg,
		ToExpPos:  c.ToExpPos,
		MinE:      c.MinE,
		MaxE:      c.MaxE,
		Crypto:    c.Crypto,
	}
}

// Option defines a functional configuration option for Context.
type Option func(*Context)

// WithPrecision sets the maximum significant digits for calculations.
func WithPrecision(p int) Option {
	return func(c *Context) {
		if p >= 1 && p <= MaxDigits {
			c.Precision = p
		}
	}
}

// WithRounding sets the rounding mode.
func WithRounding(r RoundingMode) Option {
	return func(c *Context) {
		if r <= RoundHalfFloor {
			c.Rounding = r
		}
	}
}

// WithModulo sets the modulo mode.
func WithModulo(m ModuloMode) Option {
	return func(c *Context) {
		c.Modulo = m
	}
}

// WithToExpNeg sets the negative exponential threshold.
func WithToExpNeg(e int64) Option {
	return func(c *Context) {
		if e <= 0 && e >= -ExpLimit {
			c.ToExpNeg = e
		}
	}
}

// WithToExpPos sets the positive exponential threshold.
func WithToExpPos(e int64) Option {
	return func(c *Context) {
		if e >= 0 && e <= ExpLimit {
			c.ToExpPos = e
		}
	}
}

// WithMinE sets the minimum exponent limit.
func WithMinE(e int64) Option {
	return func(c *Context) {
		if e >= -ExpLimit && e <= 0 {
			c.MinE = e
		}
	}
}

// WithMaxE sets the maximum exponent limit.
func WithMaxE(e int64) Option {
	return func(c *Context) {
		if e >= 0 && e <= ExpLimit {
			c.MaxE = e
		}
	}
}

// Config applies the provided options to the Context in-place.
func (c *Context) Config(opts ...Option) *Context {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	return c
}
