package rpn

const (
	left = iota
	right
)

type operator struct {
	priority int
	assoc    int
}

var operators = map[string]operator{
	"^": operator{4, right},
	"*": operator{3, left},
	"/": operator{3, left},
	"+": operator{2, left},
	"-": operator{2, left},
}
