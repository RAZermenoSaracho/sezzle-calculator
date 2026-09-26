package calculator

type tokenKind int

const (
	tokEOF tokenKind = iota
	tokNumber
	tokPlus
	tokMinus
	tokStar
	tokSlash
	tokCaret
	tokPercent
	tokLParen
	tokRParen
	tokSqrt
)

type token struct {
	kind tokenKind
	text string
}

var symbols = map[byte]tokenKind{
	'+': tokPlus, '-': tokMinus, '*': tokStar, '/': tokSlash,
	'^': tokCaret, '%': tokPercent, '(': tokLParen, ')': tokRParen,
}

func isDigit(c byte) bool  { return c >= '0' && c <= '9' }
func isLetter(c byte) bool { return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' }

// lex splits input into tokens. Numbers are DIGIT+ [. DIGIT+] or . DIGIT+;
// the only accepted word is the lowercase "sqrt".
func lex(input string) ([]token, error) {
	var tokens []token
	for i := 0; i < len(input); {
		c := input[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case isDigit(c) || c == '.':
			start := i
			for i < len(input) && isDigit(input[i]) {
				i++
			}
			if i < len(input) && input[i] == '.' {
				i++
				fracStart := i
				for i < len(input) && isDigit(input[i]) {
					i++
				}
				if i == fracStart {
					return nil, ErrSyntax
				}
			}
			tokens = append(tokens, token{tokNumber, input[start:i]})
		case isLetter(c):
			start := i
			for i < len(input) && isLetter(input[i]) {
				i++
			}
			if input[start:i] != "sqrt" {
				return nil, ErrSyntax
			}
			tokens = append(tokens, token{tokSqrt, "sqrt"})
		default:
			kind, ok := symbols[c]
			if !ok {
				return nil, ErrSyntax
			}
			tokens = append(tokens, token{kind, string(c)})
			i++
		}
	}
	return append(tokens, token{kind: tokEOF}), nil
}
