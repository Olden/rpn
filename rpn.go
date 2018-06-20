package rpn

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// FromInfix convert infix expressions to postfix expressions (reverse Polish notation)
func FromInfix(in string) (string, error) {
	var out []string
	buf := newStack()
	var tmp bytes.Buffer

	var prev rune
	for i, v := range in {
		if unicode.IsDigit(v) {
			tmp.WriteRune(v)
			continue
		}
		if tmp.Len() > 0 {
			out = append(out, tmp.String())
			tmp.Reset()
			prev = 0
		}
		op := string(v)
		switch op {
		case "^", "*", "/", "+", "-":
			// check unary operation
			if (op == "-" || op == "+") && (operatorsList.isOperator(string(prev)) || i == 0) {
				out = append(out, "0")
				buf.push(op)
				break
			}

			for t := buf.peak(); t != nil; t = buf.peak() {
				top, err := operatorsList.get(t.(string))
				o1, err1 := operatorsList.get(op)

				if err != nil || err1 != nil || o1.greater(top) ||
					o1.equals(top) && o1.assoc == right {
					break
				}
				out = append(out, buf.pop().(string))
			}
			buf.push(op)
			prev = v
		case "(":
			buf.push(op)
		case ")":
			for {
				if buf.top == nil {
					return "", fmt.Errorf("Invalid bracket order: %s. Not enough open bracket", in)
				}
				l := buf.pop()
				if l.(string) == "(" {
					break
				}
				out = append(out, l.(string))
			}
		}
	}

	if tmp.Len() > 0 {
		out = append(out, tmp.String())
	}

	for buf.top != nil {
		l := buf.pop()
		if l.(string) == "(" {
			return "", fmt.Errorf("Invalid bracket order: %s. Not enough closed bracket", in)
		}
		out = append(out, l.(string))
	}

	return strings.Join(out, " "), nil
}

// Calculate given postfix expression
func Calculate(in string) (float64, error) {
	if !isValidRpn(in) {
		return 0, fmt.Errorf("Invalid postfix notation: %s", in)
	}
	buf := newStack()

	for _, v := range strings.Fields(in) {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			buf.push(n)
			continue
		}

		sec := buf.pop().(float64)
		first := buf.pop().(float64)

		op, err := operatorsList.get(v)
		if err != nil {
			return 0, fmt.Errorf("Invalid postfix notation: %s", in)
		}
		buf.push(op.call(first, sec))
	}

	return buf.pop().(float64), nil
}

func isValidRpn(in string) bool {
	c := 0
	for _, v := range strings.Fields(in) {
		if _, err := strconv.ParseFloat(v, 64); err == nil {
			c++
			continue
		}

		if !operatorsList.isOperator(v) {
			continue
		}
		if c < 1 {
			return false
		}
		c--
	}

	return c == 1
}
