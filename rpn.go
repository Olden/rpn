package rpn

import (
	"math"
	"strconv"
	"strings"
)

func FromInfix(in string) string {
	out := newStack()
	buf := newStack()

	in = strings.Replace(in, " ", "", -1)

	for i := 0; i < len(in); i++ {
		op := string(in[i])

		if _, err := strconv.ParseInt(op, 10, 32); err == nil {
			out.push(op)
			continue
		}

		switch op {
		case "^", "*", "/", "+", "-":
			for {
				topChar := buf.pop()
				if topChar == nil {
					buf.push(topChar)
					break
				}
				top, ok := operators[topChar.(string)]
				if !ok {
					buf.push(topChar)
					break
				}
				if operators[op].priority < top.priority ||
					operators[op].priority == top.priority && operators[op].assoc == left {
					out.push(topChar)
				} else {
					buf.push(topChar)
					break
				}
			}
			buf.push(op)
		case "(":
			buf.push(op)
		case ")":
			for {
				l := buf.pop()
				if l.(string) != "(" {
					out.push(l)
					continue
				}
				break
			}
		}
	}

	for {
		l := buf.pop()
		if l == nil {
			break
		}
		out.push(l)
	}

	return out.string()
}

func Calculate(in string) float64 {
	buf := newStack()

	in = strings.Replace(in, " ", "", -1)

	for i := 0; i < len(in); i++ {
		op := string(in[i])

		if _, err := strconv.ParseFloat(op, 64); err == nil {
			buf.push(op)
			continue
		}

		secChar := buf.pop()
		firstChar := buf.pop()
		sec, _ := strconv.ParseFloat(secChar.(string), 64)
		first, _ := strconv.ParseFloat(firstChar.(string), 64)
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
