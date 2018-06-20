package rpn

import (
	"strings"
)

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

// Return string representation of the stack
func (n *node) string() []string {
	var out []string

	if n.prev != nil {
		return append(n.prev.string(), n.value.(string))
	}

	return append(out, n.value.(string))
}

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

// String returns string representation of stack elements separated by space
func (s *stack) string() string {
	return strings.Join(s.top.string(), " ")
}

// NewStack returns a new stack.
func newStack() *stack {
	return &stack{}
}
