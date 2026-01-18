package main

import (
	"fmt"
	"log"
	"os"

	"gihub.com/roidaradal/nlp"
	"github.com/roidaradal/fn/dict"
)

const path string = "data/1.json"

func main() {
	jsonLexer := nlp.Lexer{
		Tokens: []nlp.Token{
			{"LEFT_BRACE", `\{`},
			{"RIGHT_BRACE", `\}`},
			{"LEFT_BRACKET", `\[`},
			{"RIGHT_BRACKET", `\]`},
			{"COLON", ":"},
			{"COMMA", ","},
			{"LITERAL_TRUE", "true"},
			{"LITERAL_FALSE", "false"},
			{"LITERAL_NULL", "null"},
			{"STRING", `"[^"]*"`},
			{"NUMBER", `-?\d+(\.\d+)?`},
			{"WHITESPACE", `\s+`},
		},
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	text := string(bytes)

	replacements := dict.StringMap{
		"WHITESPACE": " ",
	}
	tokens, err := jsonLexer.Tokenize(text, replacements)
	if err != nil {
		log.Fatal(err)
	}
	for _, token := range tokens {
		fmt.Println(token)
	}
}
