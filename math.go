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
}

var operators = map[string]operator{
	"^": operator{4, right},
	"*": operator{3, left},
	"/": operator{3, left},
	"+": operator{2, left},
	"-": operator{2, left},
}

type mathFunc func(f, s float64) float64

var calcFunctions = map[string]mathFunc{
	"^": func(f, s float64) float64 {
		return math.Pow(f, s)
	},
	"*": func(f, s float64) float64 {
		return float64(f * s)
	},
	"/": func(f, s float64) float64 {
		return f / s
	},
	"+": func(f, s float64) float64 {
		return float64(f + s)
	},
	"-": func(f, s float64) float64 {
		return float64(f - s)
	},
}
