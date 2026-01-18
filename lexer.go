package nlp

import (
	"fmt"
	"regexp"

	"github.com/roidaradal/fn/dict"
)

type Token [2]string

type Lexer struct {
	Tokens []Token // [TokenType, TokenRegExp]
}

// Tokenize the given text
func (l Lexer) Tokenize(text string, replacements dict.StringMap) ([]Token, error) {
	// Prepare token patterns
	patterns := make([]*regexp.Regexp, len(l.Tokens))
	for i, token := range l.Tokens {
		pattern := "^" + token[1]
		patterns[i] = regexp.MustCompile(pattern)
	}
	tokens := make([]Token, 0)
	for len(text) > 0 {
		found := false
		for i, pattern := range patterns {
			match := pattern.FindStringIndex(text)
			if match == nil {
				continue
			}
			start, end := match[0], match[1]
			chunk := text[start:end]
			tokenType := l.Tokens[i][0]
			text = text[end:]
			if rep, ok := replacements[tokenType]; ok {
				chunk = rep
			}
			tokens = append(tokens, Token{tokenType, chunk})
			found = true
			break
		}
		if !found {
			limit := min(10, len(text))
			return nil, fmt.Errorf("failed to tokenize: %s", text[:limit])
		}
	}
	return tokens, nil
}

// Destructure token parts
func (t Token) Tuple() (string, string) {
	return t[0], t[1]
}
