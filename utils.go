package nlp

import (
	"bufio"
	"os"

	"github.com/roidaradal/fn/str"
)

// Read byte lines from given path
func ReadLineBytes(path string) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, str.WrapError("failed to open file", err)
	}
	defer file.Close()

	lines := make([][]byte, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
