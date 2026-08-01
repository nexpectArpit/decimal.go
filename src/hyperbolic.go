package decimal

// Sinh computes hyperbolic sine sinh(x) = (exp(x) - exp(-x)) / 2 using context c settings.
func (c *Context) Sinh(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() || x.d == nil {
		return &Decimal{s: 0}
	}
	if x.IsZero() {
		return &Decimal{s: x.s, e: 0, d: []int32{0}}
	}

	expX := c.Exp(x)
	negX := x.Neg()
	expNegX := c.Exp(negX)
	diff := c.Sub(expX, expNegX)
	two, _ := c.New(2)
	return c.Div(diff, two)
}

// Sinh computes hyperbolic sine using default context.
func (x *Decimal) Sinh() *Decimal {
	return globalContext.Sinh(x)
}

// HyperbolicSine is an alias for Sinh.
func (x *Decimal) HyperbolicSine() *Decimal {
	return x.Sinh()
}

// Cosh computes hyperbolic cosine cosh(x) = (exp(x) + exp(-x)) / 2 using context c settings.
func (c *Context) Cosh(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() || x.d == nil {
		return &Decimal{s: 0}
	}
	if x.IsZero() {
		return &Decimal{s: 1, e: 0, d: []int32{1}}
	}

	expX := c.Exp(x)
	negX := x.Neg()
	expNegX := c.Exp(negX)
	sum := c.Add(expX, expNegX)
	two, _ := c.New(2)
	return c.Div(sum, two)
}

// Cosh computes hyperbolic cosine using default context.
func (x *Decimal) Cosh() *Decimal {
	return globalContext.Cosh(x)
}

// HyperbolicCosine is an alias for Cosh.
func (x *Decimal) HyperbolicCosine() *Decimal {
	return x.Cosh()
}

// Tanh computes hyperbolic tangent tanh(x) = sinh(x) / cosh(x) using context c settings.
func (c *Context) Tanh(x *Decimal) *Decimal {
	sinhX := c.Sinh(x)
	coshX := c.Cosh(x)
	return c.Div(sinhX, coshX)
}

// Tanh computes hyperbolic tangent using default context.
func (x *Decimal) Tanh() *Decimal {
	return globalContext.Tanh(x)
}

// HyperbolicTangent is an alias for Tanh.
func (x *Decimal) HyperbolicTangent() *Decimal {
	return x.Tanh()
}

// Asinh computes inverse hyperbolic sine asinh(x) = ln(x + sqrt(x^2 + 1)) using context c settings.
func (c *Context) Asinh(x *Decimal) *Decimal {
	one, _ := c.New(1)
	x2 := c.Mul(x, x)
	x2Plus1 := c.Add(x2, one)
	sqrtVal := c.Sqrt(x2Plus1)
	sum := c.Add(x, sqrtVal)
	return c.Ln(sum)
}

// Asinh computes inverse hyperbolic sine using default context.
func (x *Decimal) Asinh() *Decimal {
	return globalContext.Asinh(x)
}

// Acosh computes inverse hyperbolic cosine acosh(x) = ln(x + sqrt(x^2 - 1)) using context c settings.
func (c *Context) Acosh(x *Decimal) *Decimal {
	one, _ := c.New(1)
	x2 := c.Mul(x, x)
	x2Minus1 := c.Sub(x2, one)
	sqrtVal := c.Sqrt(x2Minus1)
	sum := c.Add(x, sqrtVal)
	return c.Ln(sum)
}

// Acosh computes inverse hyperbolic cosine using default context.
func (x *Decimal) Acosh() *Decimal {
	return globalContext.Acosh(x)
}

// Atanh computes inverse hyperbolic tangent atanh(x) = 0.5 * ln((1 + x) / (1 - x)) using context c settings.
func (c *Context) Atanh(x *Decimal) *Decimal {
	one, _ := c.New(1)
	half, _ := c.New(0.5)
	num := c.Add(one, x)
	den := c.Sub(one, x)
	frac := c.Div(num, den)
	lnFrac := c.Ln(frac)
	return c.Mul(half, lnFrac)
}

// Atanh computes inverse hyperbolic tangent using default context.
func (x *Decimal) Atanh() *Decimal {
	return globalContext.Atanh(x)
}
