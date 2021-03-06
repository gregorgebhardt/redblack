package redblack

import (
	"math"
)

type Node struct {
	key         int64
	value       interface{}
	red         bool
	left, right *Node
}

func (n *Node) height() chan int {
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

func (n *Node) width() int {
	h := <-n.height()
	return int(math.Pow(2., float64(h-1)))
}

func (n *Node) min() *Node {
	if n.left != nil {
		return n.left.min()
	}
	return n
}

func (n *Node) max() *Node {
	if n.right != nil {
		return n.right.max()
	}
	return n
}

func (n *Node) isLeaf() bool {
	return n.left == nil && n.right == nil
}

func (n *Node) walkInOrder(f func(*Node)) {
	if n == nil {
		f(n)
	} else {
		n.left.walkInOrder(f)
		f(n)
		n.right.walkInOrder(f)
	}
}

func (n *Node) walkPreOrder(f func(*Node)) {
	if n == nil {
		f(n)
	} else {
		f(n)
		n.left.walkPreOrder(f)
		n.right.walkPreOrder(f)
	}
}

func (n *Node) walkPostOrder(f func(*Node)) {
	if n == nil {
		f(n)
	} else {
		n.left.walkPostOrder(f)
		n.right.walkPostOrder(f)
		f(n)
	}
}

func (n *Node) walkLevelOrder(queue []*Node, f func(*Node)) {
	f(n)
	if n != nil {
		queue = append(queue, n.left, n.right)
		queue[0].walkLevelOrder(queue[1:], f)
	}
}

func (n *Node) search(k int64) *Node {
	if n == nil {
		return nil
	} else if n.key == k {
		return n
	} else if n.key > k {
		return n.left.search(k)
	} else {
		return n.right.search(k)
	}
}

func (n *Node) searchUpper(k int64) *Node {
	if n == nil {
		return nil
	} else if n.key == k {
		return n
	} else if n.key > k {
		nc := n.left.searchUpper(k)
		if nc == nil {
			return n
		}
		return nc
	} else {
		return n.right.searchUpper(k)
	}
}

func (n *Node) searchLower(k int64) *Node {
	if n == nil {
		return nil
	} else if n.key == k {
		return n
	} else if n.key > k {
		return n.left.searchLower(k)
	} else {
		nc := n.right.searchLower(k)
		if nc == nil {
			return n
		}
		return nc
	}
}

type keyError string

func (e keyError) Error() string {
	return string(e)
}

const KeyExistsError = keyError("Key already exists in tree.")
const KeyDoesNotExistError = keyError("Key not found.")

func (n *Node) insert(key int64, value interface{}) (*Node, error) {
	if n == nil {
		return &Node{key: key, value: value, red: true}, nil
	}

	if isRed(n.left) && isRed(n.right) {
		n.flipColors()
	}

	if key == n.key {
		return nil, KeyExistsError
	} else if key < n.key {
		newNode, err := n.left.insert(key, value)
		if err != nil {
			return nil, err
		}
		n.left = newNode
	} else {
		newNode, err := n.right.insert(key, value)
		if err != nil {
			return nil, err
		}
		n.right = newNode
	}
	n = n.fixUp()

	return n, nil
}

func isRed(n *Node) bool {
	return n != nil && n.red
}

func (n *Node) flipColors() {
	n.red = !n.red
	n.left.red = !n.left.red
	n.right.red = !n.right.red
}

func (n *Node) rotateLeft() *Node {
	x := n.right
	n.right = x.left
	x.left = n
	x.red = n.red
	n.red = true
	return x
}

func (n *Node) rotateRight() *Node {
	x := n.left
	n.left = x.right
	x.right = n
	x.red = n.red
	n.red = true
	return x
}

func (n *Node) deleteMin() *Node {
	if n.left == nil {
		return nil
	}

	if !isRed(n.left) && !isRed(n.left.left) {
		n = n.moveRedLeft()
	}

	n.left = n.left.deleteMin()

	return n.fixUp()
}

func (n *Node) delete(k int64) (*Node, bool) {
	var success bool
	if k < n.key {
		if !isRed(n.left) && !isRed(n.left.left) {
			n = n.moveRedLeft()
		}
		n.left, success = n.left.delete(k)
	} else {
		if isRed(n.left) {
			n = n.rotateRight()
		}
		if k == n.key && n.right == nil {
			return nil, true
		}
		if !isRed(n.right) && !isRed(n.right.left) {
			n = n.moveRedRight()
		}
		if k == n.key {
			n.key = n.right.min().key
			n.right = n.right.deleteMin()
			success = true
		} else {
			n.right, success = n.right.delete(k)
		}
	}

	return n.fixUp(), success
}

func (n *Node) moveRedLeft() *Node {
	n.flipColors()
	if isRed(n.right.left) {
		n.right = n.right.rotateRight()
		n = n.rotateLeft()
		n.flipColors()
	}
	return n
}

func (n *Node) moveRedRight() *Node {
	n.flipColors()
	if isRed(n.left.left) {
		n = n.rotateRight()
		n.flipColors()
	}
	return n
}

func (n *Node) fixUp() *Node {
	if isRed(n.right) && !isRed(n.left) {
		n = n.rotateLeft()
	}
	if isRed(n.left) && isRed(n.left.left) {
		n = n.rotateRight()
	}
	return n
}
