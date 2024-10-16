package main

import (
	"fmt"

	"github.com/gregorgebhardt/redblack"
)

func main() {
	values := []redblack.Orderable[string]{
		redblack.Ordered("a"),
		redblack.Ordered("b"),
		redblack.Ordered("c"),
	}
	tree, err := redblack.NewTree(values, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tree)
}
