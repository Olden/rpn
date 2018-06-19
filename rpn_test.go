package rpn

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

type exp struct {
	infix  string
	rpn    string
	result float64
}

var cases = []exp{
	exp{"3 / 0", "3 0 /", math.Inf(1)},
	exp{"3 + 2", "3 2 +", 5},
	exp{"3 2 4 + +", "3 2 4 + +", 9},
	exp{" 3 + 2 ", "3 2 +", 5},
	exp{"(3 + 2)", "3 2 +", 5},
	exp{"(3 + (2))", "3 2 +", 5},
	exp{"3 + 2 * 3", "3 2 3 * +", 9},
	exp{"3 + 2 + 3 * 4", "3 2 + 3 4 * +", 17},
	exp{"3 * 2 + 3 + 4", "3 2 * 3 + 4 +", 13},
	exp{"3 + 2 * 3 + 4 * 5", "3 2 3 * + 4 5 * +", 29},
	exp{"(3 + 2) * 3", "3 2 + 3 *", 15},
	exp{"3 * 2 + 3 * (5 + 2)", "3 2 * 3 5 2 + * +", 27},
	exp{"3 * 2 + 3 * (5 + 2 * 3)", "3 2 * 3 5 2 3 * + * +", 39},
	exp{"((3 + 2) * 5)", "3 2 + 5 *", 25},
	exp{"((3 - 1) / 2) * 3 + 5", "3 1 - 2 / 3 * 5 +", 8},
	exp{"((3 - 1) / 2) * (3 + 5)", "3 1 - 2 / 3 5 + *", 8},
	exp{"((((1 * (2 + 3)) - 3) + 4) * 5)", "1 2 3 + * 3 - 4 + 5 *", 30},
	exp{"(1 + 2) * 4 + 3", "1 2 + 4 * 3 +", 15},
	exp{"3 + 4 * 2 / (1 - 5) ^ 2", "3 4 2 * 1 5 - 2 ^ / +", 3.5},
	exp{"9 ^ (1 / 2)", "9 1 2 / ^", 3},
	exp{"9 ^ ((0 - 1) / 2)", "9 0 1 - 2 / ^", 0.3333333333333333},
	exp{"((15 / (7 - (1 + 1))) * 3) - (2 + (1 + 1))", "15 7 1 1 + - / 3 * 2 1 1 + + -", 5},
}

func TestFromInfix(t *testing.T) {
	for i, e := range cases {
		r, err := FromInfix(e.infix)
		if err != nil {
			t.Error("unexpected error:", err)
		}
		if e.rpn != r {
			t.Error("case:", i,
				"\n\tinfix:       ", e.infix,
				"\n\texpected rpn:", e.rpn,
				"\n\tresult:      ", r)
		}
	}
}

func BenchmarkFromInfix(b *testing.B) {
	for j := 0; j < b.N; j++ {
		FromInfix("((15 / (7 - (1 + 1))) * 3) - (2 + (1 + 1))")
	}
}

func TestRpnCalculation(t *testing.T) {
	for i, e := range cases {
		r, err := Calculate(e.rpn)
		if err != nil {
			t.Error("unexpected error:", err)
		}
		if e.result != r {
			t.Error("case:", i,
				"\n\tinfix:          ", e.infix,
				"\n\trpn:            ", e.rpn,
				"\n\texpected result:", e.result,
				"\n\tresult:         ", r)
		}
	}
}

func BenchmarkRpnCalculation(b *testing.B) {
	for j := 0; j < b.N; j++ {
		Calculate("15 7 1 1 + - / 3 * 2 1 1 + + -")
	}
}

type invalidInfix struct {
	infix string
	err   error
}

var invalidInfixCases = []invalidInfix{
	invalidInfix{"3 + 3)", errors.New("Invalid bracket order: 3 + 3). Not enough open bracket")},
	invalidInfix{"(3 + 3", errors.New("Invalid bracket order: (3 + 3. Not enough closed bracket")},
	invalidInfix{"((((1 * (2 + 3)) - 3) + 4) * 5", errors.New("Invalid bracket order: ((((1 * (2 + 3)) - 3) + 4) * 5. Not enough closed bracket")},
}

func TestInvalidInfixInput(t *testing.T) {
	for i, e := range invalidInfixCases {
		_, err := FromInfix(e.infix)
		if e.err.Error() != err.Error() {
			t.Error("case:", i,
				"\n\tinfix:          ", e.infix,
				"\n\texpected error: ", e.err,
				"\n\tresult:         ", err)
		}
	}
}

var invalidPostfixCases = []string{
	"1 a + 3 *",
	"a 1 + 3 *",
	"+ 3",
	"2 3 ?",
	"2 3 + *",
}

func TestInvalidPostfixInput(t *testing.T) {
	for i, c := range invalidPostfixCases {
		_, err := Calculate(c)
		exp := fmt.Errorf("Invalid postfix notation: %s", c)
		if exp.Error() != err.Error() {
			t.Error("case:", i,
				"\n\trpn:            ", c,
				"\n\texpected error: ", exp,
				"\n\tresult:         ", err)
		}
	}
}

func TestPostfixValidation(t *testing.T) {
	for _, e := range invalidPostfixCases {
		if isValidRpn(e) != false {
			t.Errorf("Given rpn: %s is invalid, but it not been determined", e)
		}
	}
}

func TestNotTrimmedPostfixString(t *testing.T) {
	r, _ := Calculate(" 3 3 + ")
	if r != 6 {
		t.Error(
			"\n\trpn:             3 3 + ",
			"\n\texpected result: 6",
			"\n\tresult:          ", r)
	}
}
