package redblack

import "iter"

type Tree[V any, T Orderable[V]] struct {
	root *Node[V, T]
	num  int
}

type WalkOrder int

const (
	INORDER WalkOrder = iota
	PREORDER
	POSTORDER
	LEVELORDER
)

func (t *Tree[V, T]) Search(k V) (bool, V) {
	n := t.root.search(k)
	if n != nil {
		return true, n.Value()
	}
	return false, k
}

func (t *Tree[V, T]) SearchUpper(k V) (V, error) {
	if n := t.root.searchUpper(k); n != nil {
		return n.Value(), nil
	}
	return k, KeyDoesNotExistError
}

func (t *Tree[V, T]) SearchLower(k V) (V, error) {
	if n := t.root.searchLower(k); n != nil {
		return n.Value(), nil
	}
	return k, KeyDoesNotExistError
}

// Insert adds a new node to the tree.
// Returns true if the insertion was successful and false if the node already exists
func (t *Tree[V, T]) Insert(item T) error {
	newNode, err := t.root.insert(item)
	if err != nil {
		return err
	}
	t.num++
	t.root = newNode
	t.root.red = false
	return nil
}

func (t *Tree[V, T]) Delete(v V) (success bool) {
	if t.root == nil {
		return false
	}

	if found, _ := t.Search(v); !found {
		return false
	}

	t.root, success = t.root.delete(v)
	if t.root != nil {
		t.root.red = false
	}
	if success {
		t.num--
	}
	return
}

func (t *Tree[V, T]) DeleteMin() {
	if t.root != nil {
		t.root = t.root.deleteMin()
		t.num--
	}
}

func NewTree[V any, T Orderable[V]](items []T) (*Tree[V, T], error) {
	tree := new(Tree[V, T])
	for _, v := range items {
		if err := tree.Insert(v); err != nil {
			return nil, err
		}
	}
	return tree, nil
}

func (t *Tree[V, T]) Height() int {
	return <-t.root.height()
}

func (t *Tree[V, T]) Len() int {
	return t.num
}

func (t *Tree[V, T]) Min() interface{} {
	return t.root.min().value
}

func (t *Tree[V, T]) Max() interface{} {
	return t.root.max().value
}

func (t *Tree[V, T]) ToSortedSlice() []V {
	values := make([]V, 0, t.num)
	f := func(n *Node[V, T]) bool {
		if n != nil {
			values = append(values, n.Value())
		}
		return true
	}
	t.root.walkInOrder(f)
	return values
}

func (t *Tree[V, T]) GetTreeLevels() [][]*Node[V, T] {
	h := t.Height()
	level := make([][]*Node[V, T], h)

	thisLevel := make([]*Node[V, T], 1)
	thisLevel[0] = t.root

	for l := 0; l < h; l++ {
		level[l] = thisLevel
		nextLevel := make([]*Node[V, T], int(1<<(l+1)))
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

func (t Tree[V, T]) Sorted() iter.Seq[V] {
	return func(yield func(V) bool) {
		f := func(n *Node[V, T]) bool {
			if n != nil {
				return yield(n.Value())
			}
			return true
		}
		t.Walk(f, INORDER)
	}
}

func (t *Tree[V, T]) Walk(f func(*Node[V, T]) bool, order WalkOrder) {
	switch order {
	case INORDER:
		t.root.walkInOrder(f)
	case PREORDER:
		t.root.walkPreOrder(f)
	case POSTORDER:
		t.root.walkPostOrder(f)
	case LEVELORDER:
		t.root.walkLevelOrder(make([]*Node[V, T], 0, t.num), f)
	}
}

func (t *Tree[V, T]) checkNoRedRed() bool {
	noRedRed := true
	f := func(n *Node[V, T]) bool {
		if isRed(n) && (isRed(n.left) || isRed(n.right)) {
			noRedRed = false
		}
		return true
	}
	t.root.walkPreOrder(f)
	return noRedRed
}

func (t *Tree[V, T]) checkBlackHeight() (uint, bool) {
	blackHeightStack := make([]uint, 0, 10)
	blackHeightSame := true
	f := func(n *Node[V, T]) bool {
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
		return true
	}
	t.root.walkPostOrder(f)
	return blackHeightStack[0], blackHeightSame
}

func (t *Tree[V, T]) checkLeftLeaning() bool {
	leftLeaning := true
	f := func(n *Node[V, T]) bool {
		if n != nil && (!isRed(n.left) && isRed(n.right)) {
			leftLeaning = false
		}
		return true
	}
	t.root.walkPreOrder(f)
	return leftLeaning
}
