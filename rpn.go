package rpn

import (
	"math"
	"strconv"
	"strings"
)

// FromInfix convert infix expressions to postfix expressions (reverse Polish notation)
func FromInfix(in string) string {
	out := newStack()
	buf := newStack()

	in = strings.Replace(in, " ", "", -1)

	for i := 0; i < len(in); i++ {
		op := string(in[i])
		if _, err := strconv.ParseFloat(op, 64); err == nil {
			out.push(op)
			continue
		}

		switch op {
		case "^", "*", "/", "+", "-":
			for {
				topChar := buf.peak()
				if topChar == nil {
					break
				}
				top, isOperator := operators[topChar.(string)]
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
			for {
				l := buf.pop()
				if l.(string) == "(" {
					break
				}
				out.push(l)
			}
		}
	}

	for buf.length > 0 {
		out.push(buf.pop())
	}

	return out.string()
}

// Calculate given postfix expression
func Calculate(in string) float64 {
	buf := newStack()

	in = strings.Replace(in, " ", "", -1)

	for i := 0; i < len(in); i++ {
		op := string(in[i])

		if _, err := strconv.ParseFloat(op, 64); err == nil {
			buf.push(op)
			continue
		}

		sec, _ := strconv.ParseFloat(buf.pop().(string), 64)
		first, _ := strconv.ParseFloat(buf.pop().(string), 64)
		var r float64

		switch op {
		case "^":
			r = math.Pow(first, sec)
		case "*":
			r = float64(first * sec)
		case "/":
			r = first / sec
		case "+":
			r = float64(first + sec)
		case "-":
			r = float64(first - sec)
		}

		buf.push(strconv.FormatFloat(r, 'f', 2, 64))
	}

	res, _ := strconv.ParseFloat(buf.string(), 10)

	return res
}
