package main

import (
	"log"
	"strings"
	"unicode"
)

func divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}

func check(e error, exit bool) {
	if e != nil {
		if exit {
			log.Panic(e)
		} else {
			log.Print(e)
		}
	}
}

func removeWhitespace(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))
	for _, chr := range str {
		if !unicode.IsSpace(chr) {
			builder.WriteRune(chr)
		}
	}
	return builder.String()
}
