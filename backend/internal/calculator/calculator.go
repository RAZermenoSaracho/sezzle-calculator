// Package calculator evaluates arithmetic expressions. It knows nothing about
// HTTP or any other transport.
package calculator

import (
	"errors"
	"strings"
)

var (
	ErrEmpty          = errors.New("empty expression")
	ErrSyntax         = errors.New("invalid expression")
	ErrDivisionByZero = errors.New("division by zero")
	ErrInvalidSqrt    = errors.New("invalid square root")
	ErrInvalidNumber  = errors.New("invalid number")
	ErrTooComplex     = errors.New("expression too complex")
)

const (
	maxLength = 1000
	maxDepth  = 100
)

// Evaluate parses and evaluates expression. Grammar and precedence are
// documented in the expression-evaluation spec: postfix %, right-associative
// ^, unary sign, * /, + -.
func Evaluate(expression string) (float64, error) {
	if strings.TrimSpace(expression) == "" {
		return 0, ErrEmpty
	}
	if len(expression) > maxLength {
		return 0, ErrTooComplex
	}

	tokens, err := lex(expression)
	if err != nil {
		return 0, err
	}
	p := &parser{tokens: tokens}
	result, err := p.additive()
	if err != nil {
		return 0, err
	}
	if p.peek().kind != tokEOF {
		return 0, ErrSyntax
	}
	if result == 0 {
		result = 0 // normalize -0
	}
	return result, nil
}
