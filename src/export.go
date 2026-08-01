package decimal

// ExportDigitsToString exposes digitsToString for testing.
func ExportDigitsToString(d []int32) string {
	return digitsToString(d)
}

// ExportConvertBase exposes convertBase for testing.
func ExportConvertBase(str string, baseIn, baseOut int) []int {
	return convertBase(str, baseIn, baseOut)
}
