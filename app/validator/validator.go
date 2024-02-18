package validator

import (
	"unicode"
)

func IsValid(expression string) bool {

	for _, r := range expression {
		switch {
		case unicode.Is(unicode.Digit, r):
			continue
		case r == ' ':
			continue
		case r == '+':
			continue
		case r == '-':
			continue
		case r == '*':
			continue
		case r == '/':
			continue
		default:
			return false
		}
	}
	return true
}
