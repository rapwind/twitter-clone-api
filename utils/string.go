package utils

import (
	"strings"
)

// StringTrimFromMap get string that remove the specified characters
func StringTrimFromMap(str string, qs []string) string {
	for _, v := range qs {
		str = strings.Replace(str, v, "", -1)
	}
	return str
}

// PhoneNumberNormalization get normalization phone number from string
func PhoneNumberNormalization(phoneNumber string) string {
	// TODO
	qs := []string{" ", "+", "-", "(", ")"}
	return StringTrimFromMap(phoneNumber, qs)
}
