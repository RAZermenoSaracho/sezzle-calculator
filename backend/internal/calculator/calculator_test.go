package calculator

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"addition", "1 + 2", 3},
		{"subtraction", "5 - 8", -3},
		{"multiplication", "4 * 2.5", 10},
		{"division", "9 / 4", 2.25},
		{"multiplication before addition", "1 + 2 * 3", 7},
		{"division before subtraction", "10 - 6 / 2", 7},
		{"parentheses override precedence", "(1 + 2) * 3", 9},
		{"nested parentheses", "((1 + 2) * (3 + 4))", 21},
		{"subtraction is left associative", "10 - 4 - 3", 3},
		{"division is left associative", "100 / 10 / 5", 2},
		{"exponent", "2 ^ 3", 8},
		{"exponent is right associative", "2 ^ 3 ^ 2", 512},
		{"unary minus is below exponent", "-2 ^ 2", -4},
		{"signed exponent", "2 ^ -2", 0.25},
		{"signed exponent chain", "2 ^ -3 ^ 2", math.Pow(2, -9)},
		{"exponent binds tighter than multiplication", "2 * 3 ^ 2", 18},
		{"fractional exponent", "16 ^ 0.5", 4},
		{"unary minus", "-5 + 2", -3},
		{"double unary minus", "--5", 5},
		{"plus minus", "+-5", -5},
		{"unary before parentheses", "-(2 + 3)", -5},
		{"multiply by negative", "3 * -2", -6},
		{"decimal", "1.5 + 2.25", 3.75},
		{"leading dot decimal", ".5 + 1.25", 1.75},
		{"whitespace ignored", "  1+   2 *3 ", 7},
		{"tabs and newlines", "1\t+\n2", 3},
		{"sqrt", "sqrt(16)", 4},
		{"sqrt of expression", "sqrt(2 + 2) * 3", 6},
		{"sqrt zero", "sqrt(0)", 0},
		{"nested sqrt", "sqrt(sqrt(16))", 2},
		{"sqrt then power", "sqrt(2) ^ 2", 2},
		{"percent", "50%", 0.5},
		{"percent is not contextual", "100 + 10%", 100.1},
		{"repeated percent", "50%%", 0.005},
		{"negative percent", "-50%", -0.5},
		{"percent binds tighter than exponent base", "50% ^ 2", 0.25},
		{"percent binds tighter than exponent", "2 ^ 50%", math.Sqrt(2)},
		{"percent of parenthesized", "(1 + 1)%", 0.02},
		{"negative zero normalized", "-0", 0},
		{"reference expression", "1 + 2 - 4 * 3 * (5 - 1) ^ 0.5 + sqrt(16)", -17},
		{"sqrt with percent", "sqrt(16) + 50%", 4.5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Evaluate(tc.input)
			if err != nil {
				t.Fatalf("Evaluate(%q) error = %v", tc.input, err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Evaluate(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestEvaluateNegativeZero(t *testing.T) {
	got, err := Evaluate("-0")
	if err != nil || math.Signbit(got) {
		t.Errorf("Evaluate(\"-0\") = %v, %v; want positive zero", got, err)
	}
}

func TestEvaluateErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{"empty", "", ErrEmpty},
		{"whitespace only", "   \t", ErrEmpty},

		{"trailing operator", "1 +", ErrSyntax},
		{"leading operator", "* 2", ErrSyntax},
		{"missing closing parenthesis", "(1 + 2", ErrSyntax},
		{"extra closing parenthesis", "1 + 2)", ErrSyntax},
		{"empty parentheses", "()", ErrSyntax},
		{"two numbers", "1 2", ErrSyntax},
		{"trailing dot", "5.", ErrSyntax},
		{"double dot", "1.2.3", ErrSyntax},
		{"lone dot", ".", ErrSyntax},
		{"scientific notation", "1e3", ErrSyntax},
		{"implicit multiplication", "2(3)", ErrSyntax},
		{"unknown word", "abc", ErrSyntax},
		{"uppercase sqrt", "SQRT(4)", ErrSyntax},
		{"unsupported character", "1 $ 2", ErrSyntax},
		{"non-ASCII character", "1 × 2", ErrSyntax},
		{"sqrt without parentheses", "sqrt 16", ErrSyntax},
		{"sqrt without argument", "sqrt()", ErrSyntax},
		{"sqrt unclosed", "sqrt(16", ErrSyntax},
		{"lone percent", "%", ErrSyntax},
		{"percent before operand", "%5", ErrSyntax},
		{"double operator", "1 * * 2", ErrSyntax},
		{"missing exponent", "2 ^", ErrSyntax},

		{"division by zero", "1 / 0", ErrDivisionByZero},
		{"division by zero expression", "1 / (2 - 2)", ErrDivisionByZero},
		{"zero to negative power", "0 ^ -1", ErrDivisionByZero},

		{"sqrt of negative", "sqrt(-1)", ErrInvalidSqrt},
		{"sqrt of negative expression", "sqrt(1 - 2)", ErrInvalidSqrt},

		{"negative base fractional exponent", "(-8) ^ 0.5", ErrInvalidNumber},
		{"overflow power", "10 ^ 1000", ErrInvalidNumber},
		{"overflow literal within limit", "1" + strings.Repeat("0", 400), ErrInvalidNumber},
		{"overflow multiplication", "1" + strings.Repeat("0", 200) + " * 1" + strings.Repeat("0", 200), ErrInvalidNumber},

		{"too long", strings.Repeat("1+", 501), ErrTooComplex},
		{"too deep parentheses", strings.Repeat("(", 200) + "1" + strings.Repeat(")", 200), ErrTooComplex},
		{"too deep unary", strings.Repeat("-", 500) + "1", ErrTooComplex},
		{"far too many parentheses", strings.Repeat("(", 10000), ErrTooComplex},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Evaluate(tc.input)
			if !errors.Is(err, tc.want) {
				t.Errorf("Evaluate(%q) = %v, %v; want error %v", tc.input, got, err, tc.want)
			}
		})
	}
}

func TestEvaluateMaxNestingAllowed(t *testing.T) {
	input := strings.Repeat("(", 90) + "1" + strings.Repeat(")", 90)
	got, err := Evaluate(input)
	if err != nil || got != 1 {
		t.Errorf("Evaluate(90 nested) = %v, %v; want 1, nil", got, err)
	}
}

func FuzzEvaluate(f *testing.F) {
	for _, seed := range []string{
		"1 + 2 * 3", "(1 + 2) * 3", "-5 + 2", "2 ^ 3 ^ 2", "sqrt(16)", "50%",
		"1 / 0", "((", "))", "sqrt(", "1e999", "..", "-", "%%%", "\x00",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		got, err := Evaluate(input)
		if err == nil && (math.IsNaN(got) || math.IsInf(got, 0)) {
			t.Errorf("Evaluate(%q) returned non-finite %v without error", input, got)
		}
	})
}
