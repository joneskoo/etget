// Package htmltable implements a HTML table parser
package htmltable

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// Table represents a HTML table
type Table struct {
	Headers [][]string
	Rows    [][]string
}

// Parse parses HTML from r
func Parse(r io.Reader) (page []Table, err error) {
	n, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %s", err)
	}
	tables := getElementsByName(n, "table")
	for _, t := range tables {
		table := parseTable(t)
		page = append(page, table)
	}
	return
}

func parseTable(n *html.Node) (table Table) {
	theads := getElementsByName(n, "thead")
	if len(theads) > 0 {
		table.Headers = parseRows(theads[0])
	}
	tbodies := getElementsByName(n, "tbody")
	if len(tbodies) > 0 {
		table.Rows = parseRows(tbodies[0])
	}
	return
}

func parseRows(n *html.Node) (rows [][]string) {
	for _, tr := range getElementsByName(n, "tr") {
		elems := []string{}
		for _, td := range getElementsByName(tr, "td") {
			elems = append(elems, getTextContent(td))
		}
		rows = append(rows, elems)
	}
	return
}

func getElementsByName(n *html.Node, name string) (elements []*html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == name {
			elements = append(elements, c)
		}
		for _, el := range getElementsByName(c, name) {
			// buf := new(bytes.Buffer)
			// html.Render(buf, c)
			// fmt.Printf("name=%s, c=%s\n", name, buf)
			elements = append(elements, el)
		}
	}
	return
}

func getTextContent(n *html.Node) string {
	buf := new(bytes.Buffer)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			fmt.Fprintf(buf, c.Data)
		}
		fmt.Fprintf(buf, getTextContent(c))
	}
	return buf.String()
}
