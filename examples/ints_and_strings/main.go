package main

import (
	"fmt"

	"github.com/gregorgebhardt/redblack"
)

func main() {
	values := map[int64]string{
		1: "a",
		2: "b",
		3: "c",
	}
	tree := redblack.NewTree(values)
	fmt.Println(tree)
}
