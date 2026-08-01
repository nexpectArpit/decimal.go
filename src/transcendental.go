package decimal

// LN10String is the natural logarithm of 10 to 1025 digits matching decimal.js line 28.
const LN10String = "2.3025850929940456840179914546843642076011014886287729760333279009675726096773524802359972050895982983419677840422862486334095254650828067566662873690987816894829072083255546808437998948262331985283935053089653777326288461633662222876982198867465436674744042432743651550489343149393914796194044002221051017141748003688084012647080685567743216228355220114804663715659121373450747856947683463616792101806445070648000277502684916746550586856935673420670581136429224554405758925724208241314695689016758940256776311356919292033376587141660230105703089634572075440370847469940168269282808481184289314848524948644871927809676271275775397027668605952496716674183485704422507197965004714951050492214776567636938662976979522110718264549734772662425709429322582798502585509785265383207606726317164309505995087807523710333101197857547331541421808427543863591778117054309827482385045648019095610299291824318237525357709750539565187697510374970888692180205189339507238539205144634197265287286965110862571492198849978748873771345686209167058"

// PIString is Pi to 1025 digits matching decimal.js line 31.
const PIString = "3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679821480865132823066470938446095505822317253594081284811174502841027019385211055596446229489549303819644288109756659334461284756482337867831652712019091456485669234603486104543266482133936072602491412737245870066063155881748815209209628292540917153643678925903600113305305488204665213841469519415116094330572703657595919530921861173819326117931051185480744623799627495673518857527248912279381830119491298336733624406566430860213949463952247371907021798609437027705392171762931767523846748184676694051320005681271452635608277857713427577896091736371787214684409012249534301465495853710507922796892589235420199561121290219608640344181598136297747713099605187072113499999983729780499510597317328160963185950244594553469083026425223082533446850352619311881710100031378387528865875332083814206171776691473035982534904287554687311595628638823537875937519577818577805321712268066130019278766111959092164201989380952572010654858632789"

// getLn10 returns natural log of 10 to sd precision.
// Matches decimal.js getLn10() (lines 3156-3166).
func (c *Context) getLn10(sd int) *Decimal {
	ln10, _ := c.New(LN10String)
	return c.finalise(ln10, sd, RoundHalfUp, false)
}

// getPi returns Pi to sd precision.
// Matches decimal.js getPi() (lines 3168-3172).
func (c *Context) getPi(sd int) *Decimal {
	pi, _ := c.New(PIString)
	return c.finalise(pi, sd, RoundHalfUp, false)
}

// Ln computes natural logarithm ln(x) using context c settings.
// Matches decimal.js naturalLogarithm() / ln (lines 3398-3512).
func (c *Context) Ln(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() || x.s < 0 {
		return &Decimal{s: 0}
	}
	if x.IsZero() {
		return &Decimal{s: -1, e: 0, d: nil} // -Infinity
	}
	if x.d == nil {
		return &Decimal{s: 1, e: 0, d: nil} // +Infinity
	}
	if x.e == 0 && len(x.d) == 1 && x.d[0] == 1 {
		return &Decimal{s: 1, e: 0, d: []int32{0}} // ln(1) = 0
	}

	wpr := c.Precision + 10
	x1, _ := c.New(x)

	// Series expansion: ln((1+y)/(1-y)) = 2(y + y^3/3 + y^5/5 + ...)
	// where y = (x - 1)/(x + 1)
	one, _ := c.New(1)
	num := c.Sub(x1, one)
	den := c.Add(x1, one)
	yTerm := c.divide(num, den, wpr, RoundDown, false)
	y2 := c.finalise(c.Mul(yTerm, yTerm), wpr, RoundDown, false)

	sum, _ := c.New(yTerm)
	numerator, _ := c.New(yTerm)
	denominator := int64(3)

	for i := 0; i < 200; i++ {
		numerator = c.finalise(c.Mul(numerator, y2), wpr, RoundDown, false)
		denDec, _ := c.New(denominator)
		term := c.divide(numerator, denDec, wpr, RoundDown, false)
		sum = c.Add(sum, term)
		if len(term.d) == 0 || term.d[0] == 0 {
			break
		}
		denominator += 2
	}

	two, _ := c.New(2)
	sum = c.Mul(sum, two)

	return c.finalise(sum, c.Precision, c.Rounding, false)
}

// Ln computes natural logarithm ln(x) using default context.
func (x *Decimal) Ln() *Decimal {
	return globalContext.Ln(x)
}

// Exp computes natural exponential e^x using context c settings.
// Matches decimal.js naturalExponential() / exp (lines 3307-3396).
func (c *Context) Exp(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if x.d == nil {
		if x.s < 0 {
			return &Decimal{s: 1, e: 0, d: []int32{0}} // exp(-Inf) = 0
		}
		return &Decimal{s: 1, e: 0, d: nil} // exp(+Inf) = +Inf
	}
	if x.IsZero() {
		return &Decimal{s: 1, e: 0, d: []int32{1}} // exp(0) = 1
	}

	wpr := c.Precision + 10
	sum, _ := c.New(1)
	pow, _ := c.New(1)
	den, _ := c.New(1)

	for i := 1; i < 200; i++ {
		pow = c.finalise(c.Mul(pow, x), wpr, RoundDown, false)
		iDec, _ := c.New(i)
		den = c.Mul(den, iDec)
		term := c.divide(pow, den, wpr, RoundDown, false)
		sum = c.Add(sum, term)
		if len(term.d) == 0 || term.d[0] == 0 {
			break
		}
	}

	return c.finalise(sum, c.Precision, c.Rounding, false)
}

// Exp computes natural exponential e^x using default context.
func (x *Decimal) Exp() *Decimal {
	return globalContext.Exp(x)
}

// Log computes logarithm of x in base y using context c settings.
func (c *Context) Log(x, y *Decimal) *Decimal {
	lnX := c.Ln(x)
	lnY := c.Ln(y)
	return c.Div(lnX, lnY)
}

// Log computes logarithm of x in base y using default context.
func (x *Decimal) Log(y *Decimal) *Decimal {
	return globalContext.Log(x, y)
}
