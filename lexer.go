package nlp

import (
	"fmt"

	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/fn/io"
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
func (l Lexer) Tokenize(lines [][]byte, ignore *ds.Set[string]) ([]Token, error) {
	for _, line := range lines {
		fmt.Println(string(line))
	}
	// // Prepare token patterns
	// patterns := list.Map(l.Tokens, func(token Token) *regexp.Regexp {
	// 	return regexp.MustCompile("^" + token[1])
	// })
	// // Tokenize text by checking each pattern for match until text is fully consumed
	// tokens := make([]Token, 0)
	// for len(bytes) > 0 {
	// 	found := false
	// 	for i, pattern := range patterns {
	// 		match := pattern.FindIndex(bytes)
	// 		if match == nil {
	// 			continue // skip if not match
	// 		}
	// 		start, end := match[0], match[1]
	// 		chunk := string(bytes[start:end]) // get chunk of text matched by pattern
	// 		tokenType := l.Tokens[i][0]       // get corresponding token type
	// 		bytes = bytes[end:]               // consume chunk and get remaining text
	// 		if ignore == nil || ignore.HasNo(tokenType) {
	// 			// Add to tokens list if no ignore set or ignore set doesn't have token type
	// 			tokens = append(tokens, Token{tokenType, chunk})
	// 		}
	// 		found = true
	// 		break
	// 	}
	// 	if !found {
	// 		limit := min(10, len(bytes)) // display first 10 (or shorter) chars left of text
	// 		return nil, fmt.Errorf("failed to tokenize: %s", string(bytes[:limit]))
	// 	}
	// }
	// return tokens, nil
	return []Token{}, nil
}

// Split bytes by newline
func splitByNewline(bytes []byte) [][]byte {
	lines := make([][]byte, 0)

	return lines
}
