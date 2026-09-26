package calculator

import (
	"math"
	"strconv"
)

// parser is a recursive-descent parser that evaluates while it parses.
// Every recursion cycle in the grammar passes through unary, so depth is
// tracked there only.
type parser struct {
	tokens []token
	pos    int
	depth  int
}

func (p *parser) peek() token { return p.tokens[p.pos] }

func (p *parser) next() token {
	t := p.tokens[p.pos]
	if t.kind != tokEOF {
		p.pos++
	}
	return t
}

func finite(v float64) (float64, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, ErrInvalidNumber
	}
	return v, nil
}

// additive := term { ('+' | '-') term }
func (p *parser) additive() (float64, error) {
	left, err := p.term()
	if err != nil {
		return 0, err
	}
	for k := p.peek().kind; k == tokPlus || k == tokMinus; k = p.peek().kind {
		p.next()
		right, err := p.term()
		if err != nil {
			return 0, err
		}
		if k == tokPlus {
			left, err = finite(left + right)
		} else {
			left, err = finite(left - right)
		}
		if err != nil {
			return 0, err
		}
	}
	return left, nil
}

// term := unary { ('*' | '/') unary }
func (p *parser) term() (float64, error) {
	left, err := p.unary()
	if err != nil {
		return 0, err
	}
	for k := p.peek().kind; k == tokStar || k == tokSlash; k = p.peek().kind {
		p.next()
		right, err := p.unary()
		if err != nil {
			return 0, err
		}
		if k == tokSlash {
			if right == 0 {
				return 0, ErrDivisionByZero
			}
			left, err = finite(left / right)
		} else {
			left, err = finite(left * right)
		}
		if err != nil {
			return 0, err
		}
	}
	return left, nil
}

// unary := ('+' | '-') unary | power
func (p *parser) unary() (float64, error) {
	p.depth++
	defer func() { p.depth-- }()
	if p.depth > maxDepth {
		return 0, ErrTooComplex
	}

	switch p.peek().kind {
	case tokPlus:
		p.next()
		return p.unary()
	case tokMinus:
		p.next()
		v, err := p.unary()
		return -v, err
	}
	return p.power()
}

// power := postfix [ '^' unary ]  (right-associative; the exponent may be signed)
func (p *parser) power() (float64, error) {
	base, err := p.postfix()
	if err != nil {
		return 0, err
	}
	if p.peek().kind != tokCaret {
		return base, nil
	}
	p.next()
	exp, err := p.unary()
	if err != nil {
		return 0, err
	}
	// math.Pow(0, negative) is +Inf; report it as the division it is.
	if base == 0 && exp < 0 {
		return 0, ErrDivisionByZero
	}
	return finite(math.Pow(base, exp))
}

// postfix := primary { '%' }
func (p *parser) postfix() (float64, error) {
	v, err := p.primary()
	if err != nil {
		return 0, err
	}
	for p.peek().kind == tokPercent {
		p.next()
		v /= 100
	}
	return v, nil
}

// primary := NUMBER | '(' additive ')' | 'sqrt' '(' additive ')'
func (p *parser) primary() (float64, error) {
	switch t := p.next(); t.kind {
	case tokNumber:
		v, err := strconv.ParseFloat(t.text, 64)
		if err != nil {
			return 0, ErrInvalidNumber
		}
		return v, nil
	case tokLParen:
		return p.parenthesized()
	case tokSqrt:
		if p.next().kind != tokLParen {
			return 0, ErrSyntax
		}
		v, err := p.parenthesized()
		if err != nil {
			return 0, err
		}
		if v < 0 {
			return 0, ErrInvalidSqrt
		}
		return math.Sqrt(v), nil
	}
	return 0, ErrSyntax
}

// parenthesized parses the rest of a group after its opening parenthesis.
func (p *parser) parenthesized() (float64, error) {
	v, err := p.additive()
	if err != nil {
		return 0, err
	}
	if p.next().kind != tokRParen {
		return 0, ErrSyntax
	}
	return v, nil
}
