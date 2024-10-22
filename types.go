package redblack

import (
	"golang.org/x/exp/constraints"
)

type Orderable[T any] interface {
	// returns a value < 0 if the receiver is less than other, 0 if they are equal, and > 0 if the receiver is greater than other
	CompareTo(other T) int
	Value() T
}

type ordered[T constraints.Ordered] struct {
	value T
}

func (o ordered[T]) CompareTo(other T) int {
	if o.value == other {
		return 0
	}
	if o.value < other {
		return -1
	}
	return 1
}

func (o ordered[T]) Value() T {
	return o.value
}

func Ordered[T constraints.Ordered](value T) ordered[T] {
	return ordered[T]{value: value}
}
