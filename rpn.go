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
	out := newStack()
	buf := newStack()
	var tmp bytes.Buffer

	for _, v := range in {
		if unicode.IsDigit(v) {
			tmp.WriteRune(v)
			continue
		}
		if tmp.Len() > 0 {
			out.push(tmp.String())
			tmp.Reset()
		}

		op := string(v)
		switch op {
		case "^", "*", "/", "+", "-":
			for t := buf.peak(); t != nil; t = buf.peak() {
				top, ok := operators[t.(string)]

				if !ok || operators[op].priority > top.priority ||
					operators[op].priority == top.priority && operators[op].assoc == right {
					break
				}
				out.push(buf.pop())
			}
			buf.push(op)
		case "(":
			buf.push(op)
		case ")":
			for {
				if buf.length == 0 {
					return "", fmt.Errorf("Invalid bracket order: %s. Not enough open bracket", in)
				}
				l := buf.pop()
				if l.(string) == "(" {
					break
				}
				out.push(l)
			}
		}
	}

	if tmp.Len() > 0 {
		out.push(tmp.String())
	}

	for buf.length > 0 {
		l := buf.pop()
		if l.(string) == "(" {
			return "", fmt.Errorf("Invalid bracket order: %s. Not enough closed bracket", in)
		}
		out.push(l)
	}

	return out.string(), nil
}

// Calculate given postfix expression
func Calculate(in string) (float64, error) {
	buf := newStack()

	for _, v := range strings.Split(in, " ") {
		if len(v) == 0 {
			continue
		}

		if n, err := strconv.ParseFloat(v, 64); err == nil {
			buf.push(n)
			continue
		}

		op, ok := operators[v]
		if !ok {
			return 0, fmt.Errorf("Unknown operator: %s", v)
		}

		if buf.length < 2 {
			return 0, fmt.Errorf("Invalid postfix notation: %s", in)
		}
		sec := buf.pop().(float64)
		first := buf.pop().(float64)

		buf.push(op.call(first, sec))
	}

	return buf.pop().(float64), nil
}
