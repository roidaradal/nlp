package nlp

import (
	"fmt"

	"github.com/roidaradal/fn/io"
	"github.com/roidaradal/fn/list"
	"github.com/roidaradal/fn/str"
)

// Read cfg file and return token lines and grammar lines
func readCfgLines(path string) (tokenLines []string, grammarLines []string, err error) {
	if !io.PathExists(path) {
		return nil, nil, fmt.Errorf("path does not exist")
	}

	lines, err := io.ReadNonEmptyLines(path)
	if err != nil {
		return nil, nil, str.WrapError("failed to open file", err)
	}

	tokenLines = make([]string, 0)
	grammarLines = make([]string, 0)
	tokenMode, grammarMode := false, false
	for _, line := range lines {
		if line == "tokens:" {
			tokenMode = true
		} else if line == "grammar:" {
			grammarMode = true
		} else if grammarMode {
			grammarLines = append(grammarLines, line)
		} else if tokenMode {
			tokenLines = append(tokenLines, line)
		}
	}

	return tokenLines, grammarLines, nil
}

// Read byte lines from given path
func ReadLineBytes(path string) ([][]byte, error) {
	strLines, err := io.ReadNonEmptyLines(path)
	if err != nil {
		return nil, err
	}
	lines := list.Map(strLines, str.ToBytes)
	return lines, nil
}
