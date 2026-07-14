package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./scripts/ci dependencies|workflow|boundary|all")
		os.Exit(2)
	}
	root, err := findRepositoryRoot(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := runPolicy(os.Args[1], root); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("ci policy %s: ok\n", os.Args[1])
}
