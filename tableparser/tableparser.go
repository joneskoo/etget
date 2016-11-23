// Package tableparser parses a HTML table into rows and columns.
package tableparser

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// Header returns the header rows in the table.
func Header(r io.Reader) (rows [][]string, err error) {
	return
}

// Data returns the data rows in the table.
func Data(r io.Reader) (rows [][]string, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	parseRoot(doc)
	return
}

// parseRoot
//  -> parseTable
//    -> parseRow

func parseRoot(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "table" {
		fmt.Println("table")
		parseTable(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseRoot(c)
	}
}

func parseTable(n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "thead":
			fmt.Println("table header")
			parseTableHead(n)
		case "tbody":
			fmt.Println("table body")
			parseTableHead(n)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseTable(c)
	}
}

func parseTableHead(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "tr" {
		elems := parseRow(n, nil)
		fmt.Println("row", elems)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseTableHead(c)
	}
}

func parseRow(n *html.Node, in []string) (elems []string) {
	elems = in
	if n.Type == html.ElementNode && n.Data == "td" {
		buf := new(bytes.Buffer)
		innerText(n, buf)
		elems = append(elems, buf.String())
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		elems = parseRow(c, elems)
	}
	return
}

func innerText(n *html.Node, w io.Writer) {
	if n.Type == html.TextNode {
		fmt.Fprintf(w, n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		innerText(c, w)
	}
}
