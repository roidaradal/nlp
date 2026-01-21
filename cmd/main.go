package main

import (
	"fmt"
	"slices"
	"strings"

	"gihub.com/roidaradal/nlp"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/fn/io"
	"github.com/roidaradal/fn/list"
	"github.com/roidaradal/fn/str"
)

func main() {
	command, options := io.GetCommandOptions("")
	switch command {
	case "tokenize":
		err := cmdTokenize(options)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	default:
		fmt.Println("Usage: nlp <tokenize> (key=value)*")
	}
}

func cmdTokenize(options dict.StringMap) error {
	// Required: file={PATH} tokens={PATH}
	// Options:  ignore={TYPE1,TYPE2,...}
	tokenizeUsage := "Usage: nlp tokenize file={PATH} tokens={PATH} (ignore={TYPE1,TYPE2,...})"

	// Get paths from options
	tokenPath, filePath := "", ""
	ignore := ds.NewSet[string]()
	for k, v := range options {
		switch k {
		case "file":
			filePath = v
		case "tokens":
			tokenPath = v
		case "ignore":
			ignore.AddItems(strings.Split(v, ","))
		}
	}

	// Check if both paths are set
	if tokenPath == "" || filePath == "" {
		fmt.Println(tokenizeUsage)
		return nil
	}

	// Create lexer from tokenPath
	lexer, err := nlp.LoadLexer(tokenPath)
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
