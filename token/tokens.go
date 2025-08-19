package token

import (
	"fmt"
	"iter"
	"log"
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

//go:generate stringer -type=tokenType
type tokenType int

const (
	EOF tokenType = iota
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	STAR
)

type (
	lexemePrefix byte
	lexeme       string
)

var defaultLexemes = [...]lexeme{
	"",
	"(",
	")",
	"{",
	"}",
	",",
	".",
	"-",
	"+",
	";",
	"*",
}

type Token struct {
	tokenType tokenType
	lexeme    lexeme
	literal   string
	err       error
}

func eofToken() Token {
	return Token{tokenType: EOF, literal: "null"}
}

func unexpectedCharToken(line int, b byte) Token {
	return Token{
		err: fmt.Errorf(
			"[line %d] Error: Unexpected character: %s",
			line,
			string(b),
		),
	}
}

func (t Token) IsError() bool {
	return t.err != nil
}

func (t Token) String() string {
	if t.IsError() {
		return t.error()
	}

	return fmt.Sprintf("%s %s %s", t.tokenType, t.lexeme, t.literal)
}

func (t Token) error() string {
	return t.err.Error()
}

type (
	scanFunc    func(t *Tokenizer) Token
	lexemeIndex map[lexemePrefix]scanFunc
)

func newLexemeIndex() lexemeIndex {
	index := make(lexemeIndex, asciiStandardSize)

	for i, lexeme := range defaultLexemes[1:] {
		index[lexemePrefix(lexeme[0])] = func(t *Tokenizer) Token {
			return scanSimpleLexeme(i+1, t)
		}
	}

	return index
}

func (i lexemeIndex) find(l lexemePrefix) (f scanFunc, found bool) {
	f, found = i[l]

	return
}

const asciiStandardSize = 128

func scanSimpleLexeme(ttype int, t *Tokenizer) Token {
	return Token{
		tokenType: tokenType(ttype),
		lexeme:    defaultLexemes[ttype],
		literal:   "null",
	}
}

var mainLexemeIndex = newLexemeIndex()

func readFileContent(filename string) (content []byte, err error) {
	content, err = os.ReadFile(filename)
	if err != nil {
		err = fmt.Errorf(
			"error reading file %q: %v\n",
			filename,
			err,
		)
	}

	return
}

type Tokenizer struct {
	data    []byte
	offset  int
	lineNum int
}

func newTokenizer(filename string) (t Tokenizer, err error) {
	t.data, err = readFileContent(filename)
	t.lineNum = 1

	return
}

func (t *Tokenizer) size() int {
	return len(t.data)
}

func (t *Tokenizer) left() int {
	return t.size() - t.offset
}

func (t *Tokenizer) done() bool {
	return t.left() == 0
}

func (t Tokenizer) peek() byte {
	return t.data[t.offset]
}

func (t *Tokenizer) read() (b byte, ok bool) {
	if t.done() {
		return
	}

	b = t.data[t.offset]
	ok = true
	t.offset++

	return
}

func (t *Tokenizer) run() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		var token Token

		for {
			b, ok := t.read()

			if !ok {
				yield(eofToken())

				return
			}

			if b == '\n' {
				t.lineNum++

				continue
			}

			f, found := mainLexemeIndex.find(lexemePrefix(b))

			log.Println(found, b)

			if !found {
				token = unexpectedCharToken(t.lineNum, b)
			} else {
				token = f(t)
			}

			if !yield(token) {
				return
			}
		}
	}
}

func Tokenize(filename string) iter.Seq[Token] {
	tokenizer, err := newTokenizer(filename)
	if err != nil {
		panic(err)
	}

	return tokenizer.run()
}
