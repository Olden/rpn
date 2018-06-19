package rpn

import "math"

const (
	left assoc = iota
	right
)

type assoc int

type operator struct {
	priority int
	assoc    assoc
	call     mathFunc
}

var operators = map[string]operator{
	"^": operator{
		4,
		right,
		func(f, s float64) float64 {
			return math.Pow(f, s)
		}},
	"*": operator{
		3,
		left,
		func(f, s float64) float64 {
			return float64(f * s)
		}},
	"/": operator{
		3,
		left,
		func(f, s float64) float64 {
			return f / s
		}},
	"+": operator{
		2,
		left,
		func(f, s float64) float64 {
			return float64(f + s)
		}},
	"-": operator{
		2,
		left,
		func(f, s float64) float64 {
			return float64(f - s)
		}},
}

type mathFunc func(f, s float64) float64
