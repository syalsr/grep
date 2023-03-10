package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Usage: echo <input_text> | your_grep.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}

func matchLine(line []byte, pattern string) (bool, error) {
	str := string(line)
	switch {
	case strings.EqualFold(pattern, "\\w"):
		for _, item := range str {
			if unicode.IsDigit(item) || unicode.IsLetter(item) {
				return true, nil
			}
		}
		return false, nil
	case strings.EqualFold(pattern, "\\d"):
		for _, item := range str {
			if unicode.IsDigit(item){
				return true, nil
			}
		}
		return false, nil
	case pattern[0] == '[' && pattern[len(pattern)-1] == ']':
		return bytes.ContainsAny(line, pattern[1:len(pattern)-1]), nil
	case pattern[:2] == "[^" && pattern[len(pattern)-1] == ']':
		return !bytes.ContainsAny(line, pattern[2:len(pattern)-1]), nil
	case utf8.RuneCountInString(pattern) == 1:
		return bytes.ContainsAny(line, pattern), nil
	}
	return false, fmt.Errorf("unsupported pattern: %q", pattern)
}
