package main

import (
	"fmt"

	"github.com/gregorgebhardt/redblack/v2"
)

func main() {
	values := []redblack.Orderable[int]{
		redblack.Ordered(1),
		redblack.Ordered(2),
		redblack.Ordered(3),
	}
	tree, err := redblack.NewTree(values, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tree)
}
