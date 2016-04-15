package utils

import (
	"regexp"
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
	qs := []string{" ", "-", "(", ")"}
	phoneNumber = StringTrimFromMap(phoneNumber, qs)
	jpPhoneExtRegexp := regexp.MustCompile(`(^[\+]810)|(^0)`)
	phoneNumber = jpPhoneExtRegexp.ReplaceAllString(phoneNumber, "81")

	return phoneNumber
}
