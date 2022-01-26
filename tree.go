package redblack

import "fmt"

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

// Insert adds a new node to the tree.
// Returns true if the insertion was successful and false if the node already exists
func (t *Tree) Insert(v int) error {
	newNode, err := t.root.insert(v)
	if err != nil {
		return err
	}
	t.num++
	t.root = newNode
	t.root.red = false
	return nil
}

func (t *Tree) Delete(v int) {
	t.root = t.root.delete(v)
	t.root.red = false
}

func NewTree(values []int) *Tree {
	tree := new(Tree)
	for _, v := range values {
		if err := tree.Insert(v); err != nil {
			fmt.Errorf("Warning: error while trying to insert %v", v)
		}
	}
	return tree
}

func (t *Tree) Height() int {
	return <-t.root.height()
}

func (t *Tree) ToSortedSlice() []int {
	values := make([]int, 0, t.num)
	f := func(n *Node) {
		if n != nil {
			values = append(values, n.key)
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
