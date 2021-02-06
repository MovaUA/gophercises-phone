// Package phone provides functions to manipulate phone numbers
package phone

import "unicode"

// Norm normalizes the provided phone number
func Norm(phoneNo string) string {
	result := make([]rune, 0, len(phoneNo))
	for _, r := range phoneNo {
		if !unicode.IsDigit(r) {
			continue
		}
		result = append(result, r)
	}
	return string(result)
}
