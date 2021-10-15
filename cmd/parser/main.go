package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	err := Parse()
	if err != nil {
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

type Link struct {
	Href string
	Text string
}

func Parse() error {
	if !hasStdIn() {
		return fmt.Errorf("usage:\n $: cat example.html | go run cmd/parser/main")
	}

	doc, err := html.Parse(os.Stdin)
	if err != nil {
		return fmt.Errorf("cannot parse %w", err)
	}
	nodes := linkNodes(doc)
	for _, node := range nodes {
		fmt.Println(node)
	}

	return nil
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes((c))...)
	}
	return ret
}
