package rpn

import (
	"testing"
)

type exp struct {
	infix  string
	rpn    string
	result float64
}

var cases = []exp{
	exp{"3 + 2", "3 2 +", 5},
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
	exp{"9 ^ ((0 - 1) / 2)", "9 0 1 - 2 / ^", 0.33},
	exp{"((15 / (7 - (1 + 1))) * 3) - (2 + (1 + 1))", "15 7 1 1 + - / 3 * 2 1 1 + + -", 5},
}

func TestFromInfix(t *testing.T) {
	for i, e := range cases {
		r := FromInfix(e.infix)
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
		for i, e := range cases {
			r := FromInfix(e.infix)
			if e.rpn != r {
				b.Error("case:", i,
					"\n\tinfix:       ", e.infix,
					"\n\texpected rpn:", e.rpn,
					"\n\tresult:      ", r)
			}
		}
	}
}

func TestRpnCalculation(t *testing.T) {
	for i, e := range cases {
		r := Calculate(e.rpn)
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
		for i, e := range cases {
			r := Calculate(e.rpn)
			if e.result != r {
				b.Error("case:", i,
					"\n\tinfix:          ", e.infix,
					"\n\trpn:            ", e.rpn,
					"\n\texpected result:", e.result,
					"\n\tresult:         ", r)
			}
		}
	}
}

func TestNotTrimmedPostfixString(t *testing.T) {
	if Calculate(" 3 3 + ") != 6 {
		t.Error(
			"\n\trpn:             3 3 + ",
			"\n\texpected result: 6",
			"\n\tresult:          ", Calculate(" 3 3 + "))
	}
}
