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

//go:generate stringer -type=tokenType
type tokenType int

const (
	EOF tokenType = iota
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	STAR
	DOT
	COMMA
	PLUS
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
	"*",
	".",
	",",
	"+",
}

type Token struct {
	tokenType tokenType
	lexeme    lexeme
	literal   string
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %s", t.tokenType, t.lexeme, t.literal)
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

func (i lexemeIndex) find(l lexemePrefix) scanFunc {
	return i[l]
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
	data   []byte
	offset int
}

func newTokenizer(filename string) (t Tokenizer, err error) {
	t.data, err = readFileContent(filename)

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

var tokenEOF = Token{tokenType: EOF, literal: "null"}

func (t *Tokenizer) run() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		for {
			b, ok := t.read()

			if !ok {
				yield(tokenEOF)

				return
			}

			f := mainLexemeIndex.find(lexemePrefix(b))

			if !yield(f(t)) {
				return
			}
		}
	}
}

func Tokenize(filename string) iter.Seq[Token] {
	// return func(yield func(Token) bool) {
	tokenizer, err := newTokenizer(filename)
	if err != nil {
		panic(err)
	}

	return tokenizer.run()
	//	if len(fileContent) > 0 {
	//		panic("Scanner not implemented")
	//	} else {
	//
	//	fmt.Println("EOF  null") // Placeholder, replace this line when
	//
	// implementing the scanner
	// }
}

// }
