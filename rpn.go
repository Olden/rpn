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

	var prev rune
	for i, v := range in {
		if unicode.IsDigit(v) {
			tmp.WriteRune(v)
			continue
		}
		if tmp.Len() > 0 {
			out.push(tmp.String())
			tmp.Reset()
			prev = 0
		}
		op := string(v)
		switch op {
		case "^", "*", "/", "+", "-":
			// check unary minus
			if op == "-" && !unicode.IsDigit(prev) && prev != 0 || op == "-" && i == 0 {
				out.push("0")
				buf.push(op)
				break
			}

			for t := buf.peak(); t != nil; t = buf.peak() {
				top, ok := operators[t.(string)]

				if !ok || operators[op].priority > top.priority ||
					operators[op].priority == top.priority && operators[op].assoc == right {
					break
				}
				out.push(buf.pop())
			}
			buf.push(op)
			prev = v
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
	if !isValidRpn(in) {
		return 0, fmt.Errorf("Invalid postfix notation: %s", in)
	}
	buf := newStack()

	for _, v := range strings.Split(strings.Trim(in, " "), " ") {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			buf.push(n)
			continue
		}

		sec := buf.pop().(float64)
		first := buf.pop().(float64)

		buf.push(operators[v].call(first, sec))
	}

	return buf.pop().(float64), nil
}

func isValidRpn(in string) bool {
	c := 0
	for _, v := range strings.Split(strings.Trim(in, " "), " ") {
		if _, err := strconv.ParseFloat(v, 64); err == nil {
			c++
			continue
		}

		_, ok := operators[v]
		if ok {
			c--
			c--
			if c < 0 {
				return false
			}
			c++
		}
	}

	if c == 1 {
		return true
	}

	return false
}
