package main

import (
	"fmt"
	"os"
	"slices"

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
			ignore.AddItems(str.CommaSplit(v))
		}
	}

	// Check if both paths are set
	if tokenPath == "" || filePath == "" {
		fmt.Println("Usage: nlp tokenize file={PATH} tokens={PATH} (ignore={TYPE1,TYPE2,...})")
		return nil
	}

	// Create lexer from tokenPath
	lexer, err := nlp.LoadLexer(tokenPath)
	if err != nil {
		return err
	}

	// Open file path
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	text := string(bytes)

	// Tokenize
	tokens, err := lexer.Tokenize(text, ignore)
	if err != nil {
		return err
	}

	// Display tokens
	maxNum := len(str.Int(len(tokens)))
	maxLength := slices.Max(list.Map(tokens, func(token nlp.Token) int {
		return len(token[0])
	}))
	template := fmt.Sprintf("[%%%dd] %%-%ds : %%s\n", maxNum, maxLength)
	for i, token := range tokens {
		tokenType, chunk := token.Tuple()
		fmt.Printf(template, i+1, tokenType, chunk)
	}
	return nil
}
