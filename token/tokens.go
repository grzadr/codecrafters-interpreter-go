package token

import (
	"fmt"
	"iter"
	"os"
	"slices"
)

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
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	SLASH
	COMMENT
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
	"=",
	"==",
	"!",
	"!=",
	"<",
	"<=",
	">",
	">=",
	"/",
	"//",
}

type Token struct {
	tokenType tokenType
	lexeme    lexeme
	literal   string
	err       error
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

func newToken(ttype tokenType) Token {
	return Token{
		tokenType: ttype,
		lexeme:    defaultLexemes[int(ttype)],
		literal:   "null",
	}
}

func eofToken() Token {
	return newToken(EOF)
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

func scanIfNext(ttype tokenType, expected byte) scanFunc {
	return func(t *Tokenizer) Token {
		if next, ok := t.peek(); !ok || next != expected {
			return newToken(ttype)
		} else {
			t.skip()

			return newToken(ttype + 1)
		}
	}
}

func scanIfNextIsEqual(ttype tokenType) scanFunc {
	return scanIfNext(ttype, '=')
}

// func scanComment() scanFunc {
// 	return func(t *Tokenizer) Token {
// 		if next, ok := t.peek(); !ok || next != '/' {
// 			return newToken(SLASH)
// 		} else {
// 			return

// 		}
// 	}
// }

func scanDefault(ttype tokenType) scanFunc {
	return func(t *Tokenizer) Token {
		return newToken(ttype)
	}
}

type (
	scanFunc    func(t *Tokenizer) Token
	lexemeIndex map[lexemePrefix]scanFunc
)

func newLexemeIndex() lexemeIndex {
	index := lexemeIndex{
		'=': scanIfNextIsEqual(EQUAL),
		'!': scanIfNextIsEqual(BANG),
		'<': scanIfNextIsEqual(LESS),
		'>': scanIfNextIsEqual(GREATER),
		'/': scanIfNext(SLASH, '/'),
	}

	for i, lexeme := range defaultLexemes[1:] {
		prefix := lexemePrefix(lexeme[0])
		if _, found := index[prefix]; found {
			continue
		}

		index[prefix] = scanDefault(tokenType(i + 1))
	}

	return index
}

func (i lexemeIndex) find(l lexemePrefix) (f scanFunc, found bool) {
	f, found = i[l]

	return
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

func (t Tokenizer) size() int {
	return len(t.data)
}

func (t Tokenizer) left() int {
	return t.size() - t.offset
}

func (t Tokenizer) ok() bool {
	return t.left() > 0
}

func (t Tokenizer) current() byte {
	return t.data[t.offset]
}

func (t Tokenizer) peek() (next byte, ok bool) {
	if ok = t.ok(); ok {
		next = t.current()
	}

	return
}

func (t *Tokenizer) skip() {
	t.offset++
}

func (t *Tokenizer) read() (b byte, ok bool) {
	if b, ok = t.peek(); !ok {
		return
	}

	t.skip()

	return
}

func (t Tokenizer) rest() []byte {
	if !t.ok() {
		return []byte{}
	}

	return t.data[t.offset:]
}

func (t Tokenizer) index(b byte) int {
	return slices.Index(t.rest(), b)
}

func (t *Tokenizer) skipLine() {
	if index := t.index('\n'); index == -1 {
		t.offset = t.size()
	} else {
		t.offset += index
	}
}

func (t *Tokenizer) run() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		var token Token

	loop:
		for {
			b, ok := t.read()

			if !ok {
				yield(eofToken())

				return
			}

			switch b {
			case '\n':
				t.lineNum++

				fallthrough
			case '\t', ' ':
				continue loop
			}

			f, found := mainLexemeIndex.find(lexemePrefix(b))

			if !found {
				token = unexpectedCharToken(t.lineNum, b)
			} else {
				token = f(t)
			}

			if token.tokenType == COMMENT {
				t.skipLine()

				continue loop
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
