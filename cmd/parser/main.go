package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/html"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Printf("error %s\n", err)
		os.Exit(1)
	}
}

func hasStdIn() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	size := fi.Size()
	if size > 0 {
		return true
	}
	return false
}

func mainErr() error {
	if !hasStdIn() {
		return fmt.Errorf("usage:\n $: cat example.html | go run cmd/parser/main")
	}

	doc, err := html.Parse(os.Stdin)
	if err != nil {
		return fmt.Errorf("cannot parse %w", err)
	}
	spew.Dump(doc)
	return nil
}
