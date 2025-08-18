package token

import (
	"fmt"
	"iter"
	"os"
)

type Tokenizer struct {
	data   []byte
	offset int
}

var tokens = map[string]

type Token struct {
	class   TokenType
	lexeme  string
	literal string
	err     error
}

func Tokenize(filename string) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			yield(
				Token{
					err: fmt.Errorf(
						"Error reading file %q: %v\n",
						filename,
						err,
					),
				},
			)

			return
		}

		if len(fileContents) > 0 {
			panic("Scanner not implemented")
		} else {
			fmt.Println("EOF  null") // Placeholder, replace this line when implementing the scanner
		}
	}
}
