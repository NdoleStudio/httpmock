package requests

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"
)

type request struct{}

func (request *request) sanitizeString(value string) string {
	return strings.TrimSpace(value)
}

func (request *request) baseURL(value string) string {
	u, _ := url.Parse(value)
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}

func (request *request) isDigits(value string) bool {
	for _, c := range value {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
