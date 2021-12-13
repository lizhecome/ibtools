package util

import (
	"math/rand"
	"strings"
	"time"
)

// StringInSlice is a function similar to "x in y" Python construct
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// SpaceDelimitedStringNotGreater returns true if the first string
// is the same as the second string or does not contain any substring
// not contained in the second string (when split by space)
func SpaceDelimitedStringNotGreater(first, second string) bool {
	// Empty string is never greater
	if first == "" {
		return true
	}

	// Split the second string by space
	secondParts := strings.Split(second, " ")

	// Iterate over space delimited parts of the first string
	for _, firstPart := range strings.Split(first, " ") {
		// If the substring is not part of the second string, return false
		if !StringInSlice(firstPart, secondParts) {
			return false
		}
	}

	// The first string is the same or more restrictive
	// than the second string, return true
	return true
}

func Join(a []string, sep string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	case 2:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return a[0] + sep + a[1]
	case 3:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return a[0] + sep + a[1] + sep + a[2]
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b := make([]byte, n)
	bp := copy(b, a[0])
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}
	return string(b)
}

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
