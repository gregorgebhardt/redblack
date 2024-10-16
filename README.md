# Red-Black Tree Implementation in Go

This repository contains an implementation of a [Red-Black Tree](https://en.wikipedia.org/wiki/Red-black_tree) in Go. A Red-Black Tree is a balanced binary search tree with additional properties that ensure the tree remains balanced during insertions and deletions.

## Project Structure

- **`node.go`**: Contains the definition and methods for the tree nodes.
- **`print.go`**: Contains functions for printing the tree structure.
- **`tree.go`**: Contains the main Red-Black Tree implementation.
- **`tree_test.go`**: Contains unit tests for the Red-Black Tree implementation.
- **`examples/`**: Contains example programs that demonstrate how to use the Red-Black Tree implementation.

## Testing
To run the tests, use the go test command:

```sh
go test ./...
```

## Examples

Some examples of how the module can be used are located in the `examples` subdirectory. You can run them using the `go run` command, for example:

```sh
go run examples/ints_and_strings
```


## License
This project is licensed under the MIT License. See the LICENSE file for details.

