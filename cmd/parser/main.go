package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/html"
)

func main() {
	links, err := Parse()
	if err != nil {
		fmt.Printf("error %s\n", err)
		os.Exit(1)
	}
	spew.Dump(links)

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

// Parse will accept an html document from std in and
// return a slice of links passed from it.
func Parse() ([]Link, error) {
	if !hasStdIn() {
		return nil, fmt.Errorf("usage:\n $: cat example.html | go run cmd/parser/main")
	}

	doc, err := html.Parse(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %w", err)
	}
	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}

func text(n *html.Node) string {
	var ret string
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type != html.ElementNode {
		return ""
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = ret + text(c)
	}

	return strings.Join(strings.Fields(ret), " ")
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
