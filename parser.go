package nlp

import (
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
	Tokens   []Token  // remaining tokens
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

	// Temporary: remove whitespace for now
	// TODO: remove this step once optional whitespace added to grammars
	tokens = list.Filter(tokens, func(token Token) bool {
		return token.Type != whitespace
	})

	// Set first step as last step with tokens
	lastStep := &deriveStep{
		Sentence: sentence{p.start},
		Tokens:   tokens,
	}
	q := ds.NewQueue[*deriveStep]()
	q.Enqueue(lastStep)

	// Invariant: sentence is non-empty and first word is a variable
	// Important: replacement rule must start with terminal or is a single variable
	for q.NotEmpty() {
		step, _ := q.Dequeue()

		if len(step.Tokens) > 0 {
			lastStep = step // set as last step if has tokens
		}

		if len(step.Sentence) == 0 {
			continue // skip if empty sentence
		}

		variable := step.Sentence[0]
		for _, rule := range p.getReplacements(variable) {
			// align front and back
			equation := newEquation(rule, step.Sentence)
			result, ok := alignFront(equation, step.Tokens)
			if !ok {
				continue // skip if not aligned
			}
			emptyLeft := len(result.Sentence) == 0
			emptyRight := len(result.Tokens) == 0
			if emptyLeft && emptyRight {
				// both sides are fully consumed = success
				return nil
			} else {
				// add to queue
				q.Enqueue(result)
			}
		}

	}
	// Exit loop = queue is empty = failed to parse
	token := lastStep.Tokens[0]
	limit := min(10, len(token.Text))
	return fmt.Errorf("syntax error: unexpected %q at line %d, col %d", token.Text[:limit], token.Row+1, token.Col+1)
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

// Align sentence and tokens from the front
func alignFront(equation sentence, tokens []Token) (*deriveStep, bool) {
	limit := min(len(equation), len(tokens))
	for i := range limit {
		left, right := equation[i], tokens[i]
		if isTerminal(left) {
			// equation has terminal in front, try to match with token
			if left != right.Type {
				return nil, false
			}
		} else {
			// equation now has variable in front, stop here
			step := &deriveStep{
				Sentence: list.Copy(equation[i:]),
				Tokens:   list.Copy(tokens[i:]),
			}
			return step, true
		}
	}
	// Everything matched so far, stop here
	step := &deriveStep{
		Sentence: list.Copy(equation[limit:]),
		Tokens:   list.Copy(tokens[limit:]),
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
