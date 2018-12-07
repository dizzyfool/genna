package model

import (
	"regexp"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/jinzhu/inflection"
)

// Singular makes singular of plural english word
func Singular(input string) string {
	return inflection.Singular(input)
}

func IsUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func IsLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func ToUpper(c byte) byte {
	return c - 32
}

func ToLower(c byte) byte {
	return c + 32
}

// CamelCased converts string to camelCase
// from github.com/go-pg/pg/internal
func CamelCased(s string) string {
	r := make([]byte, 0, len(s))
	upperNext := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' {
			upperNext = true
			continue
		}
		if upperNext {
			if IsLower(c) {
				c = ToUpper(c)
			}
			upperNext = false
		}
		r = append(r, c)
	}
	return string(r)
}

// Underscore converts string to under_scored
// from github.com/go-pg/pg/internal
func Underscore(s string) string {
	r := make([]byte, 0, len(s)+5)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if IsUpper(c) {
			if i > 0 && i+1 < len(s) && (IsLower(s[i-1]) || IsLower(s[i+1])) {
				r = append(r, '_', ToLower(c))
			} else {
				r = append(r, ToLower(c))
			}
		} else {
			r = append(r, c)
		}
	}
	return string(r)
}

func ModelName(input string) string {
	splitted := camelcase.Split(CamelCased(input))

	for i, split := range splitted {
		singular := Singular(split)
		if strings.ToLower(singular) != strings.ToLower(split) {
			splitted[i] = strings.Title(singular)
			break
		}
	}

	return strings.Join(splitted, "")
}

func StructFieldName(input string) string {
	camelCased := ReplaceSuffix(CamelCased(input), "Id", "ID")

	return strings.Title(camelCased)
}

func HasUpper(input string) bool {
	for i := 0; i < len(input); i++ {
		c := input[i]
		if IsUpper(c) {
			return true
		}
	}
	return false
}

func ReplaceSuffix(input, suffix, replace string) string {
	if strings.HasSuffix(input, suffix) {
		input = input[:len(input)-len(suffix)] + replace
	}
	return input
}

func PackageName(input string) string {
	rgxp := regexp.MustCompile(`[^a-zA-Z\d]`)
	return strings.ToLower(rgxp.ReplaceAllString(input, ""))
}
