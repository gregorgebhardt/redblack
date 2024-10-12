package main

import (
	"fmt"

	"github.com/gregorgebhardt/redblack"
)

func main() {
	values := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	tree := redblack.NewTree(values)
	fmt.Println(tree)
}
