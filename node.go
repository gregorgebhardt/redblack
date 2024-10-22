package redblack

import (
	"math"
)

type Node[V any, T Orderable[V]] struct {
	value       T
	red         bool
	left, right *Node[V, T]
}

func (n *Node[V, T]) Value() V {
	return n.value.Value()
}

func (n *Node[V, T]) height() chan int {
	c := make(chan int)
	go func(c chan int) {
		if n == nil {
			c <- 0
			return
		}
		c1 := n.left.height()
		c2 := n.right.height()
		h1, h2 := <-c1, <-c2
		if h1 > h2 {
			c <- h1 + 1
		} else {
			c <- h2 + 1
		}
	}(c)
	return c
}

func (n *Node[V, T]) width() int {
	h := <-n.height()
	return int(math.Pow(2., float64(h-1)))
}

func (n *Node[V, T]) min() *Node[V, T] {
	if n.left != nil {
		return n.left.min()
	}
	return n
}

func (n *Node[V, T]) max() *Node[V, T] {
	if n.right != nil {
		return n.right.max()
	}
	return n
}

func (n *Node[V, T]) isLeaf() bool {
	return n.left == nil && n.right == nil
}

func (n *Node[V, T]) walkInOrder(f func(*Node[V, T]) bool) bool {
	if n == nil {
		return f(n)
	} else {
		return n.left.walkInOrder(f) && f(n) && n.right.walkInOrder(f)
	}
}

func (n *Node[V, T]) walkPreOrder(f func(*Node[V, T]) bool) bool {
	if n == nil {
		return f(n)
	} else {
		return f(n) && n.left.walkPreOrder(f) && n.right.walkPreOrder(f)
	}
}

func (n *Node[V, T]) walkPostOrder(f func(*Node[V, T]) bool) bool {
	if n == nil {
		return f(n)
	} else {
		return n.left.walkPostOrder(f) && n.right.walkPostOrder(f) && f(n)
	}
}

func (n *Node[V, T]) walkLevelOrder(queue []*Node[V, T], f func(*Node[V, T]) bool) bool {
	if !f(n) {
		return false
	}
	if n != nil {
		queue = append(queue, n.left, n.right)
		return queue[0].walkLevelOrder(queue[1:], f)
	}
	return true
}

func (n *Node[V, T]) search(k V) *Node[V, T] {
	if n == nil {
		return nil
	}

	if c := n.value.CompareTo(k); c == 0 {
		return n
	} else if c < 0 {
		return n.right.search(k)
	}

	return n.left.search(k)
}

func (n *Node[V, T]) searchUpper(k V) *Node[V, T] {
	if n == nil {
		return nil
	}

	if c := n.value.CompareTo(k); c == 0 {
		return n
	} else if c < 0 {
		return n.right.searchUpper(k)
	}

	nc := n.left.searchUpper(k)
	if nc == nil {
		return n
	}
	return nc
}

func (n *Node[V, T]) searchLower(k V) *Node[V, T] {
	if n == nil {
		return nil
	}

	if c := n.value.CompareTo(k); c == 0 {
		return n
	} else if c < 0 {
		nc := n.right.searchLower(k)
		if nc == nil {
			return n
		}
		return nc
	}
	return n.left.searchLower(k)
}

type keyError string

func (e keyError) Error() string {
	return string(e)
}

const (
	KeyExistsError       = keyError("Key already exists in tree.")
	KeyDoesNotExistError = keyError("Key not found.")
)

func (n *Node[V, T]) insert(item T) (*Node[V, T], error) {
	if n == nil {
		return &Node[V, T]{value: item, red: true}, nil
	}

	if isRed(n.left) && isRed(n.right) {
		n.flipColors()
	}

	if c := n.value.CompareTo(item.Value()); c == 0 {
		return nil, KeyExistsError
	} else if c < 0 {
		newNode, err := n.right.insert(item)
		if err != nil {
			return nil, err
		}
		n.right = newNode
	} else {
		newNode, err := n.left.insert(item)
		if err != nil {
			return nil, err
		}
		n.left = newNode
	}

	n = n.fixUp()

	return n, nil
}

func isRed[V any, T Orderable[V]](n *Node[V, T]) bool {
	return n != nil && n.red
}

func (n *Node[V, T]) flipColors() {
	n.red = !n.red
	n.left.red = !n.left.red
	n.right.red = !n.right.red
}

func (n *Node[V, T]) rotateLeft() *Node[V, T] {
	x := n.right
	n.right = x.left
	x.left = n
	x.red = n.red
	n.red = true
	return x
}

func (n *Node[V, T]) rotateRight() *Node[V, T] {
	x := n.left
	n.left = x.right
	x.right = n
	x.red = n.red
	n.red = true
	return x
}

func (n *Node[V, T]) deleteMin() *Node[V, T] {
	if n.left == nil {
		return nil
	}

	if !isRed(n.left) && !isRed(n.left.left) {
		n = n.moveRedLeft()
	}

	n.left = n.left.deleteMin()

	return n.fixUp()
}

func (n *Node[V, T]) delete(k V) (*Node[V, T], bool) {
	if n == nil {
		return nil, false
	}

	var success bool
	if n.value.CompareTo(k) > 0 {
		if !isRed(n.left) && !isRed(n.left.left) {
			n = n.moveRedLeft()
		}
		n.left, success = n.left.delete(k)
	} else {
		if isRed(n.left) {
			n = n.rotateRight()
		}
		if n.value.CompareTo(k) == 0 && n.right == nil {
			return nil, true
		}
		if !isRed(n.right) && n.right != nil && !isRed(n.right.left) {
			n = n.moveRedRight()
		}
		if n.value.CompareTo(k) == 0 {
			n.value = n.right.min().value
			n.right = n.right.deleteMin()
			success = true
		} else {
			n.right, success = n.right.delete(k)
		}
	}

	return n.fixUp(), success
}

func (n *Node[V, T]) moveRedLeft() *Node[V, T] {
	n.flipColors()
	if isRed(n.right.left) {
		n.right = n.right.rotateRight()
		n = n.rotateLeft()
		n.flipColors()
	}
	return n
}

func (n *Node[V, T]) moveRedRight() *Node[V, T] {
	n.flipColors()
	if isRed(n.left.left) {
		n = n.rotateRight()
		n.flipColors()
	}
	return n
}

func (n *Node[V, T]) fixUp() *Node[V, T] {
	if isRed(n.right) && !isRed(n.left) {
		n = n.rotateLeft()
	}
	if isRed(n.left) && isRed(n.left.left) {
		n = n.rotateRight()
	}
	return n
}
