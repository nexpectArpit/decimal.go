package decimal_test

import (
	"reflect"
	"testing"

	decimal "our-projectInGO/src"
)

func TestDigitsToString(t *testing.T) {
	// Single limb
	res := decimal.ExportDigitsToString([]int32{12345})
	if res != "12345" {
		t.Fatalf("expected '12345', got '%s'", res)
	}

	// Multi limb: [1, 2] -> "1" + "000000" + "2" = "10000002"
	res = decimal.ExportDigitsToString([]int32{1, 2})
	if res != "10000002" {
		t.Fatalf("expected '10000002', got '%s'", res)
	}

	// Multi limb: [1, 2000000] -> 12000000 -> trailing zeros stripped -> "12"
	res = decimal.ExportDigitsToString([]int32{1, 2000000})
	if res != "12" {
		t.Fatalf("expected '12', got '%s'", res)
	}
}

func TestConvertBase(t *testing.T) {
	// 'ff' (base 16) -> base 10 = [2, 5, 5]
	res := decimal.ExportConvertBase("ff", 16, 10)
	expected := []int{2, 5, 5}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected %v, got %v", expected, res)
	}
}
