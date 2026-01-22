package nlp

import (
	"fmt"
	"slices"
	"strings"

	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/fn/list"
	"github.com/roidaradal/fn/str"
)

// Epsilon in grammar
// const epsilon string = "EPSILON"

type Parser struct {
	*Lexer
	variables []string
	terminals []string
	rules     map[string][][]string
}

func (p Parser) Info() string {
	out := make([]string, 0)
	out = append(out, fmt.Sprintf("Terminals: %d", len(p.terminals)))
	for _, terminal := range p.terminals {
		out = append(out, fmt.Sprintf("  %s", terminal))
	}
	out = append(out, fmt.Sprintf("Variables: %d", len(p.variables)))
	for _, variable := range p.variables {
		out = append(out, fmt.Sprintf("  %s => %d", variable, len(p.rules[variable])))
		for _, rule := range p.rules[variable] {
			out = append(out, fmt.Sprintf("    => %s", strings.Join(rule, " ")))
		}
	}
	return strings.Join(out, "\n")
}

// Load Parser grammar from file
func LoadParser(path string) (*Parser, error) {
	tokenLines, grammarLines, err := readCfgLines(path)
	if err != nil {
		return nil, err
	}

	return LoadParserLines(tokenLines, grammarLines)
}

// Load Parser grammar from lines
func LoadParserLines(tokenLines, grammarLines []string) (*Parser, error) {
	lexer, err := LoadLexerLines(tokenLines)
	if err != nil {
		return nil, err
	}

	p := &Parser{
		Lexer:     lexer,
		variables: make([]string, 0),
		rules:     make(map[string][][]string),
	}
	terminals := ds.NewSet[string]()
	for _, line := range grammarLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := str.CleanSplitN(line, "=", 2)
		variable, value := parts[0], parts[1]

		p.variables = append(p.variables, variable)
		p.rules[variable] = make([][]string, 0)
		for _, part := range str.CleanSplit(value, "|") {
			rule := str.SpaceSplit(part)
			p.rules[variable] = append(p.rules[variable], rule)
			terminals.AddItems(list.Filter(rule, isTerminal))
		}
	}
	p.terminals = terminals.Items()
	slices.Sort(p.terminals)
	return p, nil
}

// Check if token is terminal value
func isTerminal(token string) bool {
	return !isVariable(token)
}

// CHeck if token is variable
func isVariable(token string) bool {
	return strings.HasPrefix(token, "<") && strings.HasSuffix(token, ">")
}

// Create new JSON parser
func NewJSONParser() (*Parser, error) {
	tokenLines := strings.Split(jsonTokens, "\n")
	grammarLines := strings.Split(jsonGrammar, "\n")
	return LoadParserLines(tokenLines, grammarLines)
}

var jsonGrammar string = `
<JSON>      =   BOOLEAN | NULL | STRING | NUMBER | <LIST> | <OBJ>
<LIST>      =   LEFT_BRACKET <JSON> <ITEMS> RIGHT_BRACKET | LEFT_BRACKET RIGHT_BRACKET
<ITEMS>     =   EPSILON | COMMA <JSON> <ITEMS>
<OBJ>       =   LEFT_BRACE STRING COLON <JSON> <ENTRIES> RIGHT_BRACE | LEFT_BRACE RIGHT_BRACE
<ENTRIES>   =   EPSILON | COMMA STRING COLON <JSON> <ENTRIES>
`
