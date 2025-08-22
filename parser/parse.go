package parser

import (
	"fmt"
	"iter"
	"log"

	"github.com/codecrafters-io/interpreter-starter-go/scanner"
)

type pullNextToken func() (scanner.Token, bool)

type expression interface {
	consume(next pullNextToken) bool
	String() string
	Err() error
}

type expressionType int

const (
	expressionTypeLiteral  expressionType = iota
	expressionTypeGrouping expressionType = iota
)

type literal struct {
	value scanner.ValueLiteral
}

func newLiteral(
	token scanner.Token,
	_ pullNextToken,
) literal {
	return literal{value: token.Value()}
}

func (l literal) String() string {
	return l.value.FmtValue()
}
func (l literal) Err() error { return nil }

func (l literal) consume(next pullNextToken) bool {
	return true
}

type grouping struct {
	expr expression
	err  error
}

func newGroping(
	token scanner.Token,
	next pullNextToken,
) (expr grouping) {
	if expr.expr = newExpression(next); expr.expr == nil {
		expr.err = fmt.Errorf("Expected expression")

		return
	}

	if expr.err = expr.expr.Err(); expr.Err() != nil {
		return
	}

	if closing, _ := next(); closing.Type() != scanner.RIGHT_BRACE {
		expr.err = fmt.Errorf("Error at '%s': Expected ')'", closing.Raw())
	}

	return
}

func (l grouping) String() string {
	return fmt.Sprintf("(group %s)", l.expr)
}
func (l grouping) Err() error { return nil }

func (l grouping) consume(next pullNextToken) bool {
	return true
}

var expressionCreators = map[scanner.TokenType]expressionType{
	scanner.TRUE:       expressionTypeLiteral,
	scanner.FALSE:      expressionTypeLiteral,
	scanner.NIL:        expressionTypeLiteral,
	scanner.NUMBER:     expressionTypeLiteral,
	scanner.STRING:     expressionTypeLiteral,
	scanner.LEFT_PAREN: expressionTypeGrouping,
}

func newExpression(next pullNextToken) expression {
	token, _ := next()
	if token.Type() == scanner.EOF {
		return nil
	}

	exprType, found := expressionCreators[token.Type()]

	if !found {
		panic(
			fmt.Sprintf(
				"token type %q does not have registered constructor",
				token.Type(),
			),
		)
	}

	switch exprType {
	case expressionTypeLiteral:
		return newLiteral(token, next)
	case expressionTypeGrouping:
		return newGroping(token, next)
	default:
		panic(fmt.Errorf("unsupported expression type %+v", exprType))
	}
}

const errCodeParsing = 65

func CmdParse(filename string) (code int, err error) {
	tokenizer, err := scanner.NewTokenizer(filename)
	if err != nil {
		return
	}

	next, stop := iter.Pull(tokenizer.Run())
	defer stop()

	for {
		expr := newExpression(next)

		if expr == nil {
			break
		}

		if err = expr.Err(); err != nil {
			code = errCodeParsing

			log.Printf("[line %d] %w", tokenizer.Line(), err)

			break
		}

		fmt.Println(expr)
	}

	return
}
