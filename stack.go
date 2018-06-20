package rpn

// Stack is a basic LIFO
type (
	stack struct {
		top *node
	}
	node struct {
		value interface{}
		prev  *node
	}
)

// Push adds an element to the stack.
func (s *stack) push(value interface{}) {
	n := &node{value, s.top}
	s.top = n
}

// Pop removes and returns an element from the stack in last to first order.
func (s *stack) pop() interface{} {
	n := s.top
	s.top = n.prev
	return n.value
}

// View the top item on the stack
func (s *stack) peak() interface{} {
	if s.top == nil {
		return nil
	}
	return s.top.value
}

// NewStack returns a new stack.
func newStack() *stack {
	return &stack{}
}
