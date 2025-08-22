package parser

import (
	"fmt"
	"iter"
	"log"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/scanner"
)

type pullNextToken func() (scanner.Token, bool)

type expression interface {
	String() string
	Err() error
}

type expressionType int

const (
	expressionTypeLiteral expressionType = iota
	expressionTypeGrouping
	expressionTypeUnary
	expressionTypeBinary
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

func (e literal) String() string {
	return e.value.FmtValue()
}
func (e literal) Err() error { return nil }

type grouping struct {
	expr expression
	err  error
}

func newGroping(
	_ scanner.Token,
	next pullNextToken,
) (expr grouping) {
	if expr.expr = newExpression(next); expr.expr == nil {
		expr.err = fmt.Errorf("Expected expression")

		return
	}

	if expr.err = expr.expr.Err(); expr.Err() != nil {
		return
	}

	closing, _ := next()

	log.Printf("closing token %+v\n", closing)

	if closing.Type() != scanner.RIGHT_PAREN {
		expr.err = fmt.Errorf("Error at '%s': Expected ')'", closing.Raw())
	}

	return
}

func (e grouping) String() string {
	return fmt.Sprintf("(group %s)", e.expr)
}

func (e grouping) Err() error { return e.err }

type unary struct {
	opr  scanner.TokenType
	expr expression
	err  error
}

func newUnary(
	token scanner.Token,
	next pullNextToken,
) (expr unary) {
	expr.opr = token.Type()
	if expr.expr = newExpression(next); expr.expr == nil {
		return
	}

	if expr.err = expr.expr.Err(); expr.Err() != nil {
		return
	}

	return
}

func (e unary) String() string {
	return fmt.Sprintf("(%s %s)", e.opr.Raw(), e.expr)
}

func (e unary) Err() error { return e.err }

type binary struct {
	opr   scanner.TokenType
	left  expression
	right expression
	err   error
}

func newBinary(
	token scanner.Token,
	next pullNextToken,
) (expr unary) {
	expr.opr = token.Type()
	if expr.expr = newExpression(next); expr.expr == nil {
		return
	}

	if expr.err = expr.expr.Err(); expr.Err() != nil {
		return
	}

	return
}

func (e binary) String() string {
	return fmt.Sprintf("(%s %s %s)", e.opr.Raw(), e.left, e.right)
}
func (e binary) Err() error { return nil }

var expressionCreators = map[scanner.TokenType]expressionType{
	scanner.TRUE:       expressionTypeLiteral,
	scanner.FALSE:      expressionTypeLiteral,
	scanner.NIL:        expressionTypeLiteral,
	scanner.NUMBER:     expressionTypeLiteral,
	scanner.STRING:     expressionTypeLiteral,
	scanner.LEFT_PAREN: expressionTypeGrouping,
	scanner.MINUS:      expressionTypeUnary,
	scanner.BANG:       expressionTypeUnary,
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
	case expressionTypeUnary:
		return newUnary(token, next)
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

		if err := expr.Err(); err != nil {
			code = errCodeParsing

			fmt.Fprintf(os.Stderr, "[line %d] %s\n", tokenizer.Line(), err)

			break
		}

		fmt.Println(expr)
	}

	return
}
