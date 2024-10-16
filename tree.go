package redblack

import "iter"

type Tree[V any, T Orderable[V]] struct {
	root *Node[V, T]
	num  int
}

// WalkOrder specifies the order in which the nodes are visited when walking the tree.
type WalkOrder int

const (
	// INORDER visits all left children (smaller keys), then the current node, then all right children (larger keys).
	// This results in keys being visited in ascending order.
	INORDER WalkOrder = iota
	// PREORDER visits the current node before its child nodes (left child before right child).
	PREORDER
	// POSTORDER visits the current node after its child nodes (left child before right child).
	POSTORDER
	// LEVELORDER visits the nodes level by level from left to right.
	LEVELORDER
)

// Search returns true if the key is found in the tree and the value of the key.
// If the key is not found, the second return value is the key itself.
func (t *Tree[V, T]) Search(k V) (bool, V) {
	n := t.root.search(k)
	if n != nil {
		return true, n.Value()
	}
	return false, k
}

// SearchUpper returns the value of the smallest key in the tree that is greater than or equal to the given key.
// Returns KeyDoesNotExistError if k > i for all i in the tree.
func (t *Tree[V, T]) SearchUpper(k V) (V, error) {
	if n := t.root.searchUpper(k); n != nil {
		return n.Value(), nil
	}
	return k, KeyDoesNotExistError
}

// SearchLower returns the value of the largest key in the tree that is less than or equal to the given key.
// Returns KeyDoesNotExistError if k < i for all i in the tree.
func (t *Tree[V, T]) SearchLower(k V) (V, error) {
	if n := t.root.searchLower(k); n != nil {
		return n.Value(), nil
	}
	return k, KeyDoesNotExistError
}

// Insert adds a new node to the tree if the item is not a duplicate of another item in the tree.
// Returns KeyExistsError if the key already exists in the tree.
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

// Delete removes a node from the tree if the key is found.
// Returns false if the key is not found.
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

// DeleteMin removes the node with the smallest key from the tree.
func (t *Tree[V, T]) DeleteMin() {
	if t.root != nil {
		t.root = t.root.deleteMin()
		t.num--
	}
}

// Creates a new red-black tree from a slice of Orderable items.
// If ignore_duplicates is true, duplicate items will be ignored otherwise a KeyExistsError will be returned.
func NewTree[V any, T Orderable[V]](items []T, ignore_duplicates bool) (*Tree[V, T], error) {
	tree := new(Tree[V, T])
	for _, v := range items {
		if err := tree.Insert(v); err != nil {
			if err == KeyExistsError && ignore_duplicates {
				continue
			}
			return nil, err
		}
	}
	return tree, nil
}

// Height return the height of the tree.
// The height of a tree is the number of edges on the longest path between the root and a leaf.
func (t *Tree[V, T]) Height() int {
	return <-t.root.height()
}

// Len returns the number of nodes in the tree.
func (t *Tree[V, T]) Len() int {
	return t.num
}

// Min returns the smallest key in the tree.
func (t *Tree[V, T]) Min() V {
	return t.root.min().Value()
}

// Max returns the largest key in the tree.
func (t *Tree[V, T]) Max() V {
	return t.root.max().Value()
}

// Returns a sorted slice of the keys in the tree.
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

// Returns each level of the tree as a slice of nodes.
// Ordered from root to leaves.
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

// Returns an iterator that yields the keys in the tree in sorted order.
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

// Walks the tree in the specified order and calls the given function for each node.
// If the function returns false, the walk is stopped.
// The order can be INORDER, PREORDER, POSTORDER or LEVELORDER.
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
