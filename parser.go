package nlp

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/fn/list"
	"github.com/roidaradal/fn/str"
)

const (
	epsilon    string = "EPSILON"
	whitespace string = "WHITESPACE"
)

type sentence = []string

type Parser struct {
	*Lexer
	variables []string
	terminals []string
	rules     map[string][]sentence
	start     string
}

type deriveStep struct {
	Sentence sentence // combination of variables and terminals
	Tokens   sentence // remaining tokens
}

func (s deriveStep) String() string {
	left := strings.Join(s.Sentence, " ")
	right := strings.Join(s.Tokens, " ")
	return fmt.Sprintf("%s | %s", left, right)
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
		terminals: make([]string, 0),
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
	if len(p.variables) > 0 {
		p.start = p.variables[0] // start = first variable
	}
	if terminals.Len() > 0 {
		p.terminals = terminals.Items()
		slices.Sort(p.terminals)
	}
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

// Tokenize given list of lines, and parse the resulting tokens
func (p *Parser) Parse(lines [][]byte, ignore *ds.Set[string]) error {
	// Tokenize the lines
	tokens, err := p.Lexer.Tokenize(lines, ignore)
	if err != nil {
		return err
	}

	// Parse the tokens
	items := list.Map(tokens, Token.GetType)
	// Temporary: remove whitespace for now
	// TODO: remove this step once optional whitespace added to grammars
	items = list.Filter(items, func(item string) bool {
		return item != whitespace
	})

	q := ds.NewQueue[*deriveStep]()
	q.Enqueue(&deriveStep{
		Sentence: sentence{p.start},
		Tokens:   items,
	})
	errParse := errors.New("parse error")

	// Invariant: sentence is non-empty and first word is a variable
	// Important: replacement rule must start with terminal or is a single variable
	for q.NotEmpty() {
		step, _ := q.Dequeue()
		variable := step.Sentence[0]
		for _, rule := range p.getReplacements(variable) {
			result, ok := align(newEquation(rule, step.Sentence), step.Tokens)
			if !ok {
				continue // skip if not aligned
			}
			emptyLeft := len(result.Sentence) == 0
			emptyRight := len(result.Tokens) == 0
			if emptyLeft && emptyRight {
				// both sides are fully consumed = success
				return nil
			} else if !emptyLeft && !emptyRight {
				// both sides are not empty = add to queue
				q.Enqueue(result)
			}
		}
	}
	return errParse
}

// Create new equation
func newEquation(rule, prev sentence) sentence {
	equation := make(sentence, 0)
	equation = append(equation, rule...)     // add replacement rule to front
	equation = append(equation, prev[1:]...) // add rest of sentence to end
	equation = list.Filter(equation, func(token string) bool {
		return token != epsilon // filter out epsilon
	})
	return equation
}

// Align sentence and tokens
func align(equation, tokens sentence) (*deriveStep, bool) {
	limit := min(len(equation), len(tokens))
	for i := range limit {
		left, right := equation[i], tokens[i]
		if isTerminal(left) {
			// equation has terminal in front, try to match with token
			if left != right {
				return nil, false
			}
		} else {
			// equation now has variable in front, stop here
			step := &deriveStep{
				Sentence: equation[i:],
				Tokens:   tokens[i:],
			}
			return step, true
		}
	}
	step := &deriveStep{
		Sentence: equation[limit:],
		Tokens:   tokens[limit:],
	}
	return step, true
}

// Get the replacement rules for given variable
func (p *Parser) getReplacements(variable string) []sentence {
	q := ds.NewQueue[string]()
	q.Enqueue(variable)

	replacements := make([]sentence, 0)
	for q.NotEmpty() {
		variable, _ = q.Dequeue()
		for _, rule := range p.rules[variable] {
			first := rule[0]
			if isTerminal(first) {
				// Add to replacement if rule's first word is terminal
				replacements = append(replacements, rule)
			} else {
				// If variable, enqueue so that it can be expanded
				q.Enqueue(first)
			}
		}
	}
	return replacements
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
