package conf

const floatZero float64 = 1e-6

// floatEqual 判断相等
func floatEqual(a, b float64) bool {
	if a > b {
		return a-b < floatZero
	} else {
		return b-a < floatZero
	}
}

// floatLessEqual 判断小于等于
func floatLessEqual(a, b float64) bool {
	if a < b {
		return true
	} else {
		return a-b < floatZero
	}
}

// floatGreaterEqual 判断大于等于
func floatGreaterEqual(a, b float64) bool {
	if a > b {
		return true
	} else {
		return b-a < floatZero
	}
}
