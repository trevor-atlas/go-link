package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents a link (<a href="...">) in an HTML
// document.
type Link struct {
	Href string
	Text string
}

// Parse : parse a html document into a slice of Links
func Parse(reader io.Reader) ([]Link, error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func buildLink(node *html.Node) Link {
	var result Link
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			result.Href = attr.Val
			break
		}
	}
	result.Text = text(node)

	return result
}

func text(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	if node.Type != html.TextNode {
		return ""
	}
	var result string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		result += text(child)
	}
	return strings.Join(strings.Fields(result), " ")
}

func linkNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}
	var result []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		result = append(result, linkNodes(child)...)
	}
	return result
}
