package utils

import (
	"crypto/rand"
	"fmt"
	"html"
	"strings"
	"time"
	"unicode/utf8"
)

// EscapeParam trime string and escapes special characters
func EscapeParam(str string) string {
	return html.EscapeString(strings.Trim(str, " \t\n\r"))
}

// GetPartIn split string by space and return parts limits by len
func GetPartIn(str string, len int) string {
	result := strings.Fields(str)
	re := ""
	for _, _s := range result {
		re = fmt.Sprintf("%s %s", re, _s)
		if utf8.RuneCountInString(re) > len {
			return strings.Trim(re, " \t\n\r")
		}
	}
	return re
}

// GenTUID generate new custom UUID:len12
func GenTUID() string {
	buff := make([]byte, 2)
	_, err := rand.Read(buff)
	if err == nil {
		return fmt.Sprintf("%x%s", buff, fmt.Sprintf("%x", uint64(time.Now().UTC().UnixNano()))[8:])
	}

	return fmt.Sprintf("%x", uint64(time.Now().UTC().UnixNano()))[4:]
}
