package token

import (
	"fmt"
	"iter"
	"os"
)

// type TokenType interface {
// 	Name() string
// 	Lexeme() string
// }

// type TokenType struct {
// 	name   string
// 	lexeme string
// }

// var LEFT_PAREN = TokenType{name: "LEFT_PAREN", lexeme: "("}

// const (
// 	LEFT_PAREN  = "("
// 	RIGHT_PAREN = ")"
// 	EOF         = ""
// )

type TokenType int

const (
	EOF TokenType = iota
	LEFT_PAREN
	RIGHT_PAREN
)

var tokenLexemes = [...]string{
	"",
	"(",
	")",
}

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   string
	err       error
}

type Tokenizer struct {
	data   []byte
	offset int
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
