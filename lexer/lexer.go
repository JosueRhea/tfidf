package lexer

import (
	"strings"
)

type Tokenizer interface {
	Tokenize(text string) []string
}

type EnTokenizer struct {
}

func (s *EnTokenizer) Tokenize(text string) []string {
	// Split the text into words using whitespace as the delimiter
	words := strings.Fields(text)

	// Remove any punctuation from each word
	for i := 0; i < len(words); i++ {
		words[i] = strings.TrimFunc(words[i], func(r rune) bool {
			return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'))
		})
	}

	// Remove any empty words
	result := make([]string, 0, len(words))
	for _, word := range words {
		if word != "" {
			result = append(result, word)
		}
	}

	return result
}
