package redblack

import (
	"fmt"
)

type Tree struct {
	root *Node
	num  int
}

type WalkOrder int

const (
	INORDER WalkOrder = iota
	PREORDER
	POSTORDER
	LEVELORDER
)

func (t *Tree) Search(k int64) bool {
	return t.root.search(k) != nil
}

func (t *Tree) SearchUpper(k int64) (int64, error) {
	if n := t.root.searchUpper(k); n != nil {
		return n.key, nil
	}
	return 0, KeyDoesNotExistError
}

func (t *Tree) SearchLower(k int64) (int64, error) {
	if n := t.root.searchLower(k); n != nil {
		return n.key, nil
	}
	return 0, KeyDoesNotExistError
}

// Insert adds a new node to the tree.
// Returns true if the insertion was successful and false if the node already exists
func (t *Tree) Insert(key int64, value interface{}) error {
	newNode, err := t.root.insert(key, value)
	if err != nil {
		return err
	}
	t.num++
	t.root = newNode
	t.root.red = false
	return nil
}

func (t *Tree) Delete(v int64) (success bool) {
	t.root, success = t.root.delete(v)
	t.root.red = false
	if success {
		t.num--
	}
	return success
}

func (t *Tree) DeleteMin() {
	if t.root != nil {
		t.root.deleteMin()
		t.num--
	}
}

func NewTree(items map[int64]interface{}) *Tree {
	tree := new(Tree)
	for k, v := range items {
		if err := tree.Insert(k, v); err != nil {
			fmt.Errorf("Warning: error while trying to insert %v", v)
		}
	}
	return tree
}

func (t *Tree) Height() int {
	return <-t.root.height()
}

func (t *Tree) Len() int {
	return t.num
}

func (t *Tree) ToSortedSlice() []interface{} {
	values := make([]interface{}, 0, t.num)
	f := func(n *Node) {
		if n != nil {
			values = append(values, n.value)
		}
	}
	t.root.walkInOrder(f)
	return values
}

func (t *Tree) GetTreeLevels() [][]*Node {
	h := t.Height()
	level := make([][]*Node, h)

	thisLevel := make([]*Node, 1)
	thisLevel[0] = t.root

	for l := 0; l < h; l++ {
		level[l] = thisLevel
		nextLevel := make([]*Node, int(1<<(l+1)))
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

func (t *Tree) Walk(f func(*Node), order WalkOrder) {
	switch order {
	case INORDER:
		t.root.walkInOrder(f)
	case PREORDER:
		t.root.walkPreOrder(f)
	case POSTORDER:
		t.root.walkPostOrder(f)
	case LEVELORDER:
		t.root.walkLevelOrder(make([]*Node, 0, t.num), f)
	}
}

func (t *Tree) checkRedRed() bool {
	noRedRed := true
	f := func(n *Node) {
		if isRed(n) && (isRed(n.left) || isRed(n.right)) {
			noRedRed = false
		}
	}
	t.root.walkPreOrder(f)
	return noRedRed
}

func (t *Tree) checkBlackHeight() (uint, bool) {
	blackHeightStack := make([]uint, 0, 10)
	blackHeightSame := true
	f := func(n *Node) {
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

func (t *Tree) checkLeftLeaning() bool {
	leftLeaning := true
	f := func(n *Node) {
		if n != nil && (!isRed(n.left) && isRed(n.right)) {
			leftLeaning = false
		}
	}
	t.root.walkPreOrder(f)
	return leftLeaning
}
