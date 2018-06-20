package rpn

import (
	"fmt"
	"math"
)

const (
	left assoc = iota
	right
)

type assoc int

type mathFunc func(f, s float64) float64

type operator struct {
	priority int
	assoc    assoc
	call     mathFunc
}

func (o *operator) greater(o2 *operator) bool {
	return o.priority > o2.priority
}

func (o *operator) equals(o2 *operator) bool {
	return o.priority == o2.priority
}

type operators map[string]*operator

func (o operators) isOperator(in string) bool {
	_, ok := o[in]

	return ok
}

func (o operators) get(in string) (*operator, error) {
	if !o.isOperator(in) {
		return nil, fmt.Errorf("Unknown operator: %s", in)
	}

	return o[in], nil
}

var operatorsList = operators{
	"^": &operator{
		3,
		right,
		func(f, s float64) float64 {
			return math.Pow(f, s)
		}},
	"*": &operator{
		2,
		left,
		func(f, s float64) float64 {
			return f * s
		}},
	"/": &operator{
		2,
		left,
		func(f, s float64) float64 {
			return f / s
		}},
	"+": &operator{
		1,
		left,
		func(f, s float64) float64 {
			return f + s
		}},
	"-": &operator{
		1,
		left,
		func(f, s float64) float64 {
			return f - s
		}},
}
