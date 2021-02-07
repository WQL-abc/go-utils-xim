package xim

import (
	"fmt"
	"strings"
)

type bigram struct {
	a, b rune
}

// String - to string
func (b *bigram) String() string {
	return fmt.Sprintf("%c%c", b.a, b.b)
}

// Biunigrams - returns bigram and unigram tokens from s.
func Biunigrams(s string) []string {
	tokens := make([]string, 0, 32)

	for bigram := range toBigrams(s) {
		tokens = append(tokens, fmt.Sprintf("%c%c", bigram.a, bigram.b))
	}

	for unigram := range toUnigrams(s) {
		tokens = append(tokens, fmt.Sprintf("%c", unigram))
	}

	return tokens
}

// Bigrams returns bigram tokens from s.
func Bigrams(s string) []string {
	tokens := make([]string, 0, 32)

	for bigram := range toBigrams(s) {
		tokens = append(tokens, fmt.Sprintf("%c%c", bigram.a, bigram.b))
	}

	return tokens
}

// Prefixes - returns prefix tokens from s.
func Prefixes(s string) []string {
	return tokenize(s, false)
}

// Suffixes - returns suffix tokens from s.
func Suffixes(s string) []string {
	return tokenize(s, true)
}

func tokenize(s string, isSuffix bool) []string {
	tokenMap := make(map[string]struct{})
	runes := make([]rune, 0, 64)
	for _, w := range strings.Split(s, " ") {
		if w == "" {
			continue
		}

		if isSuffix {
			w = reverse(w)
		}

		runes = runes[0:0]

		for _, c := range w {
			runes = append(runes, c)
			tokenMap[string(runes)] = struct{}{}
		}
	}

	tokens := make([]string, 0, 32)

	for suf := range tokenMap {
		tokens = append(tokens, suf)
	}

	return tokens
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func toBigrams(value string) map[bigram]struct{} {
	result := make(map[bigram]struct{})
	var prev rune
	for i, r := range value {
		if i > 0 && prev != ' ' && r != ' ' {
			result[bigram{prev, r}] = struct{}{}
		}
		prev = r
	}
	return result
}

func toUnigrams(value string) map[rune]struct{} {
	result := make(map[rune]struct{})
	for _, r := range value {
		if r == ' ' {
			continue
		}
		result[r] = struct{}{}
	}
	return result
}
