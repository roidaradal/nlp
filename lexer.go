package nlp

import (
	"fmt"
	"regexp"

	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/fn/io"
	"github.com/roidaradal/fn/list"
	"github.com/roidaradal/fn/str"
)

type Token [2]string

// Destructure token parts
func (t Token) Tuple() (string, string) {
	return t[0], t[1]
}

type Lexer struct {
	Tokens []Token // [TokenType, TokenRegExp]
}

// Load Lexer token types from file
func LoadLexer(path string) (*Lexer, error) {
	if !io.PathExists(path) {
		return nil, fmt.Errorf("path does not exist")
	}

	lines, err := io.ReadNonEmptyLines(path)
	if err != nil {
		return nil, str.WrapError("failed to open file", err)
	}

	tokens := make([]Token, 0)
	for _, line := range lines {
		parts := str.CleanSplitN(line, ":", 2)
		tokens = append(tokens, Token{parts[0], parts[1]})
	}

	lexer := &Lexer{Tokens: tokens}
	return lexer, nil
}

// Tokenize the given text
func (l Lexer) Tokenize(text string, ignore *ds.Set[string]) ([]Token, error) {
	// Prepare token patterns
	patterns := list.Map(l.Tokens, func(token Token) *regexp.Regexp {
		return regexp.MustCompile("^" + token[1])
	})
	// Tokenize text by checking each pattern for match until text is fully consumed
	tokens := make([]Token, 0)
	for len(text) > 0 {
		found := false
		for i, pattern := range patterns {
			match := pattern.FindStringIndex(text)
			if match == nil {
				continue // skip if not match
			}
			start, end := match[0], match[1]
			chunk := text[start:end]    // get chunk of text matched by pattern
			tokenType := l.Tokens[i][0] // get corresponding token type
			text = text[end:]           // consume chunk and get remaining textf
			if ignore == nil || ignore.HasNo(tokenType) {
				// Add to tokens list if no ignore set or ignore set doesn't have token type
				tokens = append(tokens, Token{tokenType, chunk})
			}
			found = true
			break
		}
		if !found {
			limit := min(10, len(text)) // display first 10 (or shorter) chars left of text
			return nil, fmt.Errorf("failed to tokenize: %s", text[:limit])
		}
	}
	return tokens, nil
}
