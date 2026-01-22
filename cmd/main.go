package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/fn/io"
	"github.com/roidaradal/fn/list"
	"github.com/roidaradal/fn/str"
	"github.com/roidaradal/nlp"
)

func main() {
	var err error
	command, options := io.GetCommandOptions("")
	switch command {
	case "tokenize":
		err = cmdTokenize(options)
	case "parse":
		err = cmdParse(options)
	default:
		fmt.Println("Usage: nlp <tokenize|parse> (key=value)*")
	}
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

// Tokenize command handler
func cmdTokenize(options dict.StringMap) error {
	// Required: file={PATH} cfg={PATH}
	// Options:  ignore={TYPE1,TYPE2,...}
	tokenizeUsage := "Usage: nlp tokenize file={PATH} cfg={PATH} (ignore={TYPE1,TYPE2,...})"

	// Get paths from options
	cfgPath, filePath := "", ""
	ignore := ds.NewSet[string]()
	for k, v := range options {
		switch k {
		case "file":
			filePath = v
		case "cfg":
			cfgPath = v
		case "ignore":
			ignore.AddItems(strings.Split(v, ","))
		}
	}

	// Check if both paths are set
	if cfgPath == "" || filePath == "" {
		fmt.Println(tokenizeUsage)
		return nil
	}

	// Create lexer from cfgPath
	lexer, err := nlp.LoadLexer(cfgPath)
	if err != nil {
		return err
	}

	// Open file path
	lines, err := nlp.ReadLineBytes(filePath)
	if err != nil {
		return err
	}

	// Tokenize
	tokens, err := lexer.Tokenize(lines, ignore)
	if err != nil {
		return err
	}
	numTokens := len(tokens)
	fmt.Println("Tokens:", numTokens)

	if numTokens == 0 {
		return nil
	}

	// Display tokens
	maxNum := len(str.Int(len(tokens)))
	maxLength := slices.Max(list.Map(tokens, func(token nlp.Token) int {
		return len(token.Type)
	}))
	maxCoords := slices.Max(list.Map(tokens, func(token nlp.Token) int {
		return len(token.Coords())
	}))
	template := fmt.Sprintf("[%%%dd] %%-%ds : %%-%ds %%s\n", maxNum, maxLength, maxCoords)
	for i, token := range tokens {
		fmt.Printf(template, i+1, token.Type, token.Coords(), token.Text)
	}
	return nil
}

// Parse command handler
func cmdParse(options dict.StringMap) error {
	// Required: file={PATH} cfg={PATH}
	parseUsage := "Usage: nlp parse file={PATH} cfg={PATH}"

	// Get paths from options
	cfgPath, filePath := "", ""
	for k, v := range options {
		switch k {
		case "file":
			filePath = v
		case "cfg":
			cfgPath = v
		}
	}

	// Check if both parts are set
	if cfgPath == "" || filePath == "" {
		fmt.Println(parseUsage)
		return nil
	}

	// Create parser from cfgPath
	parser, err := nlp.LoadParser(cfgPath)
	if err != nil {
		return err
	}

	fmt.Println(parser.Info())

	return nil
}
