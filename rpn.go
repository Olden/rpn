package rpn

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"unicode"
)

// FromInfix convert infix expressions to postfix expressions (reverse Polish notation)
func FromInfix(in string) string {
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
				top, isOperator := operators[t.(string)]
				if !isOperator {
					break
				}
				if operators[op].priority < top.priority ||
					operators[op].priority == top.priority && operators[op].assoc == left {
					out.push(buf.pop())
				} else {
					break
				}
			}
			buf.push(op)
		case "(":
			buf.push(op)
		case ")":
			for l := buf.pop(); l.(string) != "("; l = buf.pop() {
				out.push(l)
			}
		}
	}

	if tmp.Len() > 0 {
		out.push(tmp.String())
	}

	for buf.length > 0 {
		out.push(buf.pop())
	}

	return out.string()
}

// Calculate given postfix expression
func Calculate(in string) float64 {
	buf := newStack()

	for _, v := range strings.Split(in, " ") {
		if len(v) == 0 {
			continue
		}

		if _, err := strconv.ParseFloat(v, 64); err == nil {
			buf.push(v)
			continue
		}

		sec, _ := strconv.ParseFloat(buf.pop().(string), 64)
		first, _ := strconv.ParseFloat(buf.pop().(string), 64)

		calc, isOperator := calcFunctions[v]
		if !isOperator {
			log.Fatalf("Unknown operator: %s", v)
		}

		buf.push(strconv.FormatFloat(calc(first, sec), 'f', 2, 64))
	}

	res, _ := strconv.ParseFloat(buf.string(), 10)

	return res
}
