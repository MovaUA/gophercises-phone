// Package phone provides functions to manipulate phone numbers
package phone

import (
	"bytes"
	"unicode"
)

// Norm normalizes the provided phone number
func Norm(phoneNo string) string {
	b := bytes.NewBuffer(make([]byte, 0, len(phoneNo)))
	for _, r := range phoneNo {
		if !unicode.IsDigit(r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
