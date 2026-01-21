package nlp

import (
	"fmt"
	"regexp"

	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/fn/io"
	"github.com/roidaradal/fn/list"
	"github.com/roidaradal/fn/str"
)

type TokenType [2]string // [TokenType, TokenRegExp]

// Destructure token type parts
func (t TokenType) Tuple() (string, string) {
	return t[0], t[1]
}

type Token struct {
	Type string
	Text string
	Row  int
	Col  int
}

// String containing (row, col)
func (t Token) Coords() string {
	return fmt.Sprintf("(%d, %d)", t.Row+1, t.Col+1)
}

type Lexer struct {
	TokenTypes []TokenType
	patterns   []*regexp.Regexp
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

	types := make([]TokenType, 0)
	for _, line := range lines {
		parts := str.CleanSplitN(line, ":", 2)
		types = append(types, TokenType{parts[0], parts[1]})
	}

	lexer := &Lexer{TokenTypes: types}
	return lexer, nil
}

// Tokenize the given list of lines
func (l *Lexer) Tokenize(lines [][]byte, ignore *ds.Set[string]) ([]Token, error) {
	// Prepare token patterns
	l.patterns = list.Map(l.TokenTypes, func(pair TokenType) *regexp.Regexp {
		return regexp.MustCompile("^" + pair[1])
	})

	tokens := make([]Token, 0)
	for row, line := range lines {
		lineTokens, err := l.tokenizeLine(row, line, ignore)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, lineTokens...)
	}
	return tokens, nil
}

// Tokenize one line
func (l *Lexer) tokenizeLine(row int, line []byte, ignore *ds.Set[string]) ([]Token, error) {
	// Tokenize line by checking each pattern for match until line is fully consumed
	tokens := make([]Token, 0)
	col := 0
	for len(line) > 0 {
		found := false
		for i, pattern := range l.patterns {
			match := pattern.FindIndex(line)
			if match == nil {
				continue // skip if not match
			}
			start, end := match[0], match[1]
			chunk := string(line[start:end]) // get chunk of text matched by pattern
			tokenType := l.TokenTypes[i][0]  // get corresponding token type
			line = line[end:]                // consume chunk and get remaining line
			if ignore == nil || ignore.HasNo(tokenType) {
				// Add to tokens list if no ignore set or ignore set doesnt have token type
				tokens = append(tokens, Token{Type: tokenType, Text: chunk, Row: row, Col: col})
			}
			col += end - start
			found = true
			break
		}
		if !found {
			limit := min(10, len(line))
			return nil, fmt.Errorf("failed to tokenize at line %d, col %d: %s", row+1, col+1, string(line[:limit]))
		}
	}
	return tokens, nil
}
