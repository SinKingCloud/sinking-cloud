package commerce

import "strconv"

// GetBusinessCode 工商效验码补全
func GetBusinessCode(code string, all bool) string {
	if len(code) != 14 {
		return ""
	}
	p := 10
	for i := 0; i < 14; i++ {
		num, _ := strconv.Atoi(string(code[i]))
		numerator := (p + num) % 10
		if numerator == 0 {
			numerator = 10
		}
		p = (numerator * 2) % 11
	}
	p = 11 - p
	if all {
		return code + strconv.Itoa(p)
	}
	return strconv.Itoa(p)
}
