package rpn

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"
)

var cases = []struct {
	infix  string
	rpn    string
	result float64
}{
	{"3 / 0", "3 0 /", math.Inf(1)},
	{"3 + 2", "3 2 +", 5},
	{"3 + -2", "3 0 2 - +", 1},
	{"3 - -2", "3 0 2 - -", 5},
	{"-3 + -2", "0 3 - 0 2 - +", -5},
	{"-3 + +2", "0 3 - 0 2 + +", -1},
	{"-3 + -0 + -2", "0 3 - 0 0 - + 0 2 - +", -5},
	{"-(3 + -2)*-1", "0 3 0 2 - + 0 1 - * -", 1},
	{"3 - - - - 5", "3 0 0 0 5 - - - -", 8},
	{" 3 + 2 ", "3 2 +", 5},
	{"(3 + 2)", "3 2 +", 5},
	{"(3 + (2))", "3 2 +", 5},
	{"3 + 2 * 3", "3 2 3 * +", 9},
	{"3 + 2 + 3 * 4", "3 2 + 3 4 * +", 17},
	{"3 * 2 + 3 + 4", "3 2 * 3 + 4 +", 13},
	{"3 + 2 * 3 + 4 * 5", "3 2 3 * + 4 5 * +", 29},
	{"(3 + 2) * 3", "3 2 + 3 *", 15},
	{"3 * 2 + 3 * (5 + 2)", "3 2 * 3 5 2 + * +", 27},
	{"3 * 2 + 3 * (5 + 2 * 3)", "3 2 * 3 5 2 3 * + * +", 39},
	{"((3 + 2) * 5)", "3 2 + 5 *", 25},
	{"((3 - 1) / 2) * 3 + 5", "3 1 - 2 / 3 * 5 +", 8},
	{"((3 - 1) / 2) * (3 + 5)", "3 1 - 2 / 3 5 + *", 8},
	{"((((1 * (2 + 3)) - 3) + 4) * 5)", "1 2 3 + * 3 - 4 + 5 *", 30},
	{"(1 + 2) * 4 + 3", "1 2 + 4 * 3 +", 15},
	{"3 + 4 * 2 / (1 - 5) ^ 2", "3 4 2 * 1 5 - 2 ^ / +", 3.5},
	{"9 ^ (1 / 2)", "9 1 2 / ^", 3},
	{"9 ^ ((0 - 1) / 2)", "9 0 1 - 2 / ^", 0.3333333333333333},
	{"((15 / (7 - (1 + 1))) * 3) - (2 + (1 + 1))", "15 7 1 1 + - / 3 * 2 1 1 + + -", 5},
}

func TestFromInfix(t *testing.T) {
	for i, e := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r, err := FromInfix(e.infix)
			if err != nil {
				t.Error("unexpected error:", err)
			}
			if e.rpn != r {
				t.Error("\n\tinfix:       ", e.infix,
					"\n\texpected rpn:", e.rpn,
					"\n\tresult:      ", r)
			}
		})
	}
}

func BenchmarkFromInfix(b *testing.B) {
	for j := 0; j < b.N; j++ {
		FromInfix("((15 / (7 - (1 + 1))) * 3) - (2 + (1 + 1))")
	}
}

func TestRpnCalculation(t *testing.T) {
	for i, e := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r, err := Calculate(e.rpn)
			if err != nil {
				t.Error("unexpected error:", err)
			}
			if e.result != r {
				t.Error("\n\tinfix:          ", e.infix,
					"\n\trpn:            ", e.rpn,
					"\n\texpected result:", e.result,
					"\n\tresult:         ", r)
			}
		})
	}
}

func BenchmarkRpnCalculation(b *testing.B) {
	for j := 0; j < b.N; j++ {
		Calculate("15 7 1 1 + - / 3 * 2 1 1 + + -")
	}
}

var invalidInfixCases = []struct {
	infix string
	err   error
}{
	{"3 + 3)", errors.New("Invalid bracket order: 3 + 3). Not enough open bracket")},
	{"(3 + 3", errors.New("Invalid bracket order: (3 + 3. Not enough closed bracket")},
	{"((((1 * (2 + 3)) - 3) + 4) * 5", errors.New("Invalid bracket order: ((((1 * (2 + 3)) - 3) + 4) * 5. Not enough closed bracket")},
}

func TestInvalidInfixInput(t *testing.T) {
	for i, e := range invalidInfixCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := FromInfix(e.infix)
			if e.err.Error() != err.Error() {
				t.Error("\n\tinfix:          ", e.infix,
					"\n\texpected error: ", e.err,
					"\n\tresult:         ", err)
			}
		})
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
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := Calculate(c)
			exp := fmt.Errorf("Invalid postfix notation: %s", c)
			if err == nil || exp.Error() != err.Error() {
				t.Error("\n\trpn:            ", c,
					"\n\texpected error: ", exp,
					"\n\tresult:         ", err)
			}
		})
	}
}

func TestPostfixValidation(t *testing.T) {
	for i, e := range invalidPostfixCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if isValidRpn(e) != false {
				t.Errorf("Given rpn: %s is invalid, but it not been determined", e)
			}
		})
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
