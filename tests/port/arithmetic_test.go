package decimal_test

import (
	"testing"

	decimal "our-projectInGO/src"
)

func TestAddition(t *testing.T) {
	x, _ := decimal.New("123.456")
	y, _ := decimal.New("789.123")
	res := x.Add(y)

	if res.Sign() != 1 || res.Exponent() != 2 {
		t.Fatalf("expected addition result sign=1, e=2, got sign=%d, e=%d", res.Sign(), res.Exponent())
	}
	if len(res.Coefficients()) < 2 || res.Coefficients()[0] != 912 || res.Coefficients()[1] != 5790000 {
		t.Fatalf("expected limbs [912, 5790000], got %v", res.Coefficients())
	}
}

func TestSubtraction(t *testing.T) {
	x, _ := decimal.New("1000")
	y, _ := decimal.New("1")
	res := x.Sub(y)

	if res.Sign() != 1 || res.Exponent() != 2 {
		t.Fatalf("expected subtraction result sign=1 e=2, got sign=%d e=%d", res.Sign(), res.Exponent())
	}
}

func TestMultiplication(t *testing.T) {
	x, _ := decimal.New("12.34")
	y, _ := decimal.New("5.678")
	res := x.Mul(y)

	// 12.34 * 5.678 = 70.06652
	if res.Sign() != 1 || res.Exponent() != 1 {
		t.Fatalf("expected mul sign=1 e=1, got sign=%d e=%d", res.Sign(), res.Exponent())
	}
}

func TestDivision(t *testing.T) {
	x, _ := decimal.New("355")
	y, _ := decimal.New("113")
	res := x.Div(y)

	if res.Sign() != 1 || res.Exponent() != 0 {
		t.Fatalf("expected div sign=1 e=0, got sign=%d e=%d", res.Sign(), res.Exponent())
	}
}

func TestComparison(t *testing.T) {
	x, _ := decimal.New("100")
	y, _ := decimal.New("99.999")

	if !x.Gt(y) {
		t.Fatalf("expected 100 > 99.999")
	}
	if y.Gte(x) {
		t.Fatalf("expected 99.999 < 100")
	}
	if !x.Eq(x) {
		t.Fatalf("expected 100 == 100")
	}
}
