package decimal

import (
	"math"
)

const (
	// PI constant from decimal.js with 1000 digits of precision
	PI = "3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679821480865132823066470938446095505822317253594081284811174502841027019385211055596446229489549303819644288109756659334461284756482337867831652712019091456485669234603486104543266482133936072602491412737245870066063155881748815209209628292540917153643678925903600113305305488204665213841469519415116094330572703657595919530921861173819326117931051185480744623799627495673518857527248912279381830119491298336733624406566430860213949463952247371907021798609437027705392171762931767523846748184676694051320005681271452635608277857713427577896091736371787214684409012249534301465495853710507922796892589235420199561121290219608640344181598136297747713099605187072113499999983729780499510597317328160963185950244594553469083026425223082533446850352619311881710100031378387528865875332083814206171776691473035982534904287554687311595628638823537875937519577818577805321712268066130019278766111959092164201989380952572010654858632789"
	
	// LN10 constant from decimal.js with 1000 digits of precision
	LN10 = "2.3025850929940456840179914546843642076011014886287729760333279009675726096773524802359972050895982983419677840422862486334095254650828067566662873690987816894829072083255546808437998948262331985283935053089653777326288461633662222876982198867465436674744042432743651550489343149393914796194044002221051017141748003688084012647080685567743216228355220114804663715659121373450747856947683463616792101806445070648000277502684916746550586856935673420670581136429224554405758925724208241314695689016758940256776311356919292033376587141660230105703089634572075440370847469940168269282808481184289314848524948644871927809676271275775397027668605952496716674183485704422507197965004714951050492214776567636938662976979522110718264549734772662425709429322582798502585509785265383207606726317164309505995087807523710333101197857547331541421808427543863591778117054309827482385045648019095610299291824318237525357709750539565187697510374970888692180205189339507238539205144634197265287286965110862571492198849978748873771345686209167058"
)

// getPi computes pi to sd digits (cached or computed).
func (c *Context) getPi(sd int) *Decimal {
	pi, _ := c.New(PI)
	return c.finalise(pi, sd, RoundHalfUp, false)
}

// getLn10 computes ln(10) to sd digits.
func (c *Context) getLn10(sd int) *Decimal {
	ln10, _ := c.New(LN10)
	return c.finalise(ln10, sd, RoundHalfUp, false)
}

// Ln computes natural logarithm ln(x) using context c settings.
// Matches decimal.js naturalLogarithm() / ln (lines 3398-3512).
func (c *Context) Ln(x *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0}
	}
	if x.IsZero() {
		return &Decimal{s: -1, e: 0, d: nil} // -Infinity
	}
	if x.s < 0 {
		return &Decimal{s: 0} // ln(-positive) = NaN
	}
	if x.d == nil {
		return &Decimal{s: 1, e: 0, d: nil} // +Infinity
	}
	if x.e == 0 && len(x.d) == 1 && x.d[0] == 1 {
		return &Decimal{s: 1, e: 0, d: []int32{0}} // ln(1) = 0
	}

	wpr := c.Precision + 80
	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundHalfUp

	y, _ := evalCtx.New(x)
	n := int64(1)

	// Argument reduction matching decimal.js lines 3434-3451
	digits := digitsToString(y.d)
	c0 := digits[0]
	c1 := byte('0')
	if len(digits) > 1 {
		c1 = digits[1]
	}

	for (c0 < '7' && c0 != '1') || (c0 == '1' && c1 > '3') {
		if y.IsZero() {
			break
		}
		y = evalCtx.Mul(y, x)
		if n > 1000 {
			break
		}
		digits = digitsToString(y.d)
		c0 = digits[0]
		if len(digits) > 1 {
			c1 = digits[1]
		} else {
			c1 = '0'
		}
		n++
	}

	eVal := y.e
	var normStr string
	if c0 > '1' {
		normStr = "0." + digits
		eVal++
	} else {
		normStr = string(c0) + "." + digits[1:]
	}

	xNorm, _ := evalCtx.New(normStr)

	// Taylor series: ln(xNorm) = ln((1+w)/(1-w)) = 2(w + w^3/3 + w^5/5 + ...)
	one, _ := evalCtx.New(1)
	num := evalCtx.Sub(xNorm, one)
	den := evalCtx.Add(xNorm, one)
	wTerm := evalCtx.divide(num, den, wpr, RoundDown, false)
	w2 := evalCtx.finalise(evalCtx.Mul(wTerm, wTerm), wpr, RoundDown, false)

	sum, _ := evalCtx.New(wTerm)
	numerator, _ := evalCtx.New(wTerm)
	denominator := int64(3)

	for i := 0; i < 300; i++ {
		numerator = evalCtx.finalise(evalCtx.Mul(numerator, w2), wpr, RoundDown, false)
		denDec, _ := evalCtx.New(denominator)
		term := evalCtx.divide(numerator, denDec, wpr, RoundDown, false)
		nextSum := evalCtx.Add(sum, term)
		if digitsToString(nextSum.d)[:minInt(wpr, len(digitsToString(nextSum.d)))] ==
			digitsToString(sum.d)[:minInt(wpr, len(digitsToString(sum.d)))] {
			sum = nextSum
			break
		}
		sum = nextSum
		denominator += 2
	}

	two, _ := evalCtx.New(2)
	sum = evalCtx.Mul(sum, two)

	// Reverse argument reduction matching decimal.js lines 3496-3507
	if eVal != 0 {
		ln10 := evalCtx.getLn10(wpr)
		eDec, _ := evalCtx.New(eVal)
		sum = evalCtx.Add(sum, evalCtx.Mul(eDec, ln10))
	}

	nDec, _ := evalCtx.New(n)
	sum = evalCtx.divide(sum, nDec, wpr, RoundDown, false)

	return c.finalise(sum, c.Precision, c.Rounding, true)
}

// Ln computes natural logarithm ln(x) using default context.
func (x *Decimal) Ln() *Decimal {
	return globalContext.Ln(x)
}

// Exp computes natural exponential e^x using context c settings.
// Matches decimal.js naturalExponential() / exp (lines 3307-3370).
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
	if x.e < -int64(c.Precision+2) {
		if x.s < 0 && (c.Rounding == RoundDown || c.Rounding == RoundFloor) {
			nines := make([]int32, (c.Precision+6)/7)
			for i := range nines {
				nines[i] = Base - 1
			}
			res := &Decimal{s: 1, e: -1, d: nines}
			return c.finalise(res, c.Precision, c.Rounding, true)
		}
		return &Decimal{s: 1, e: 0, d: []int32{1}}
	}
	if x.e >= 17 {
		if x.s < 0 {
			return &Decimal{s: 1, e: 0, d: []int32{0}} // exp(-large) = 0
		}
		return &Decimal{s: 1, e: 0, d: nil} // exp(+large) = +Inf
	}

	wpr := c.Precision + 60
	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	xWork, _ := evalCtx.New(x)
	tFactor, _ := evalCtx.New(0.03125) // 1 / 32 = 1 / 2^5
	k := 0

	// Argument reduction matching decimal.js lines 3334-3339
	for xWork.e > -2 {
		if xWork.IsZero() {
			break
		}
		xWork = evalCtx.Mul(xWork, tFactor)
		k += 5
		if k > 5000 {
			break
		}
	}

	sum, _ := evalCtx.New(1)
	pow, _ := evalCtx.New(1)
	den, _ := evalCtx.New(1)
	iVal := int64(0)

	for {
		iVal++
		pow = evalCtx.finalise(evalCtx.Mul(pow, xWork), wpr, RoundDown, false)
		iDec, _ := evalCtx.New(iVal)
		den = evalCtx.Mul(den, iDec)
		term := evalCtx.divide(pow, den, wpr, RoundDown, false)
		nextSum := evalCtx.Add(sum, term)

		if digitsToString(nextSum.d)[:minInt(wpr, len(digitsToString(nextSum.d)))] ==
			digitsToString(sum.d)[:minInt(wpr, len(digitsToString(sum.d)))] {
			sum = nextSum
			break
		}
		sum = nextSum
		if iVal > 300 {
			break
		}
	}

	// Reverse argument reduction: square k times matching decimal.js line 3355
	for j := 0; j < k; j++ {
		sum = evalCtx.finalise(evalCtx.Mul(sum, sum), wpr, RoundDown, false)
	}

	return c.finalise(sum, c.Precision, c.Rounding, true)
}

// Exp computes natural exponential e^x using default context.
func (x *Decimal) Exp() *Decimal {
	return globalContext.Exp(x)
}

// Log computes logarithm of x in base y using context c settings.
// Matches decimal.js logarithm / log (lines 1131-1185).
func (c *Context) Log(x, y *Decimal) *Decimal {
	if x == nil || x.IsNaN() {
		return &Decimal{s: 0} // NaN
	}
	if x.IsZero() {
		return &Decimal{s: -1, e: 0, d: nil} // -Infinity
	}
	if x.s < 0 {
		return &Decimal{s: 0} // NaN
	}

	// Default base is 10 if y is nil
	if y == nil {
		y, _ = c.New(10)
	}

	if y.IsNaN() || y.s < 0 || y.IsZero() {
		return &Decimal{s: 0} // NaN
	}

	if y.e == 0 && len(y.d) == 1 && y.d[0] == 1 {
		return &Decimal{s: 0} // NaN if base == 1
	}
	if !x.IsFinite() {
		return &Decimal{s: 1, e: 0, d: nil} // +Infinity
	}
	one, _ := c.New(1)
	if x.Eq(one) {
		return &Decimal{s: 1, e: 0, d: []int32{0}} // log(1) = 0
	}

	wpr := c.Precision + 35
	evalCtx := c.Clone()
	evalCtx.Precision = wpr
	evalCtx.Rounding = RoundDown

	lnX := evalCtx.Ln(x)

	ten, _ := evalCtx.New(10)
	if y.Eq(ten) && len(x.d) == 1 && x.d[0] == 1 {
		kDec, _ := c.New(x.e)
		return c.finalise(kDec, c.Precision, c.Rounding, false)
	}
	var lnY *Decimal
	if y.Eq(ten) {
		lnY = evalCtx.getLn10(wpr + 10)
	} else {
		lnY = evalCtx.Ln(y)
	}

	res := evalCtx.divide(lnX, lnY, wpr, RoundDown, false)

	// Check if x is an exact power y^(num/den) for small integer ratio num/den
	if resF, err := res.Float64(); err == nil && !math.IsNaN(resF) && !math.IsInf(resF, 0) {
		for den := int64(1); den <= 40; den++ {
			num := int64(math.Round(resF * float64(den)))
			if num != 0 && mathAbs(num) <= 100 {
				yPow := c.intPow(y, num, c.Precision+12)
				xPow := c.intPow(x, den, c.Precision+12)
				if yPow.Eq(xPow) {
					numDec, _ := c.New(num)
					denDec, _ := c.New(den)
					fracDec := c.divide(numDec, denDec, c.Precision+20, RoundDown, false)
					return c.finalise(fracDec, c.Precision, c.Rounding, false)
				}
			}
		}
	}

	return c.finalise(res, c.Precision, c.Rounding, true)
}

// Log computes logarithm of x in base y using default context.
func (x *Decimal) Log(y *Decimal) *Decimal {
	return globalContext.Log(x, y)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
