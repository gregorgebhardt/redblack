package redblack

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Tree[Key constraints.Ordered, T any] struct {
	root *Node[Key, T]
	num  int
}

type WalkOrder int

const (
	INORDER WalkOrder = iota
	PREORDER
	POSTORDER
	LEVELORDER
)

func (t *Tree[Key, T]) Search(k Key) bool {
	return t.root.search(k) != nil
}

func (t *Tree[Key, T]) SearchUpper(k Key) (Key, error) {
	if n := t.root.searchUpper(k); n != nil {
		return n.key, nil
	}
	return k, KeyDoesNotExistError
}

func (t *Tree[Key, T]) SearchLower(k Key) (Key, error) {
	if n := t.root.searchLower(k); n != nil {
		return n.key, nil
	}
	return k, KeyDoesNotExistError
}

// Insert adds a new node to the tree.
// Returns true if the insertion was successful and false if the node already exists
func (t *Tree[Key, T]) Insert(key Key, value T) error {
	newNode, err := t.root.insert(key, value)
	if err != nil {
		return err
	}
	t.num++
	t.root = newNode
	t.root.red = false
	return nil
}

func (t *Tree[Key, T]) Delete(v Key) (success bool) {
	t.root, success = t.root.delete(v)
	t.root.red = false
	if success {
		t.num--
	}
	return success
}

func (t *Tree[Key, T]) DeleteMin() {
	if t.root != nil {
		t.root.deleteMin()
		t.num--
	}
}

func NewTree[Key constraints.Ordered, T any](items map[Key]T) *Tree[Key, T] {
	tree := new(Tree[Key, T])
	for k, v := range items {
		if err := tree.Insert(k, v); err != nil {
			fmt.Errorf("Warning: error while trying to insert %v", v)
		}
	}
	return tree
}

func (t *Tree[Key, T]) Height() int {
	return <-t.root.height()
}

func (t *Tree[Key, T]) Len() int {
	return t.num
}

func (t *Tree[Key, T]) Min() interface{} {
	return t.root.min().value
}

func (t *Tree[Key, T]) Max() interface{} {
	return t.root.max().value
}

func (t *Tree[Key, T]) ToSortedSlice() []interface{} {
	values := make([]interface{}, 0, t.num)
	f := func(n *Node[Key, T]) {
		if n != nil {
			values = append(values, n.value)
		}
	}
	t.root.walkInOrder(f)
	return values
}

func (t *Tree[Key, T]) GetTreeLevels() [][]*Node[Key, T] {
	h := t.Height()
	level := make([][]*Node[Key, T], h)

	thisLevel := make([]*Node[Key, T], 1)
	thisLevel[0] = t.root

	for l := 0; l < h; l++ {
		level[l] = thisLevel
		nextLevel := make([]*Node[Key, T], int(1<<(l+1)))
		for i, v := range thisLevel {
			if v == nil {
				continue
			}
			iNext := 2 * i
			nextLevel[iNext] = v.left
			nextLevel[iNext+1] = v.right
		}
		thisLevel = nextLevel
	}

	return level
}

func (t *Tree[Key, T]) Walk(f func(*Node[Key, T]), order WalkOrder) {
	switch order {
	case INORDER:
		t.root.walkInOrder(f)
	case PREORDER:
		t.root.walkPreOrder(f)
	case POSTORDER:
		t.root.walkPostOrder(f)
	case LEVELORDER:
		t.root.walkLevelOrder(make([]*Node[Key, T], 0, t.num), f)
	}
}

func (t *Tree[Key, T]) checkRedRed() bool {
	noRedRed := true
	f := func(n *Node[Key, T]) {
		if isRed(n) && (isRed(n.left) || isRed(n.right)) {
			noRedRed = false
		}
	}
	t.root.walkPreOrder(f)
	return noRedRed
}

func (t *Tree[Key, T]) checkBlackHeight() (uint, bool) {
	blackHeightStack := make([]uint, 0, 10)
	blackHeightSame := true
	f := func(n *Node[Key, T]) {
		if n == nil {
			blackHeightStack = append(blackHeightStack, 1)
		} else {
			// check the heights of both siblings
			last := len(blackHeightStack) - 1
			if blackHeightStack[last] != blackHeightStack[last-1] {
				blackHeightSame = false
			}
			if n.red {
				blackHeightStack = blackHeightStack[:last]
			} else {
				blackHeightStack = blackHeightStack[:last]
				blackHeightStack[last-1]++
			}
		}
	}
	t.root.walkPostOrder(f)
	return blackHeightStack[0], blackHeightSame
}

func (t *Tree[Key, T]) checkLeftLeaning() bool {
	leftLeaning := true
	f := func(n *Node[Key, T]) {
		if n != nil && (!isRed(n.left) && isRed(n.right)) {
			leftLeaning = false
		}
	}
	t.root.walkPreOrder(f)
	return leftLeaning
}
