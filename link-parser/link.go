package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

var res []Link

func Parse(file io.Reader) ([]Link, error) {

	doc, err := html.Parse(file)
	if err != nil {
		return nil, err
	}

	getAtag(doc)

	return res, nil
}

func getTextInElement(n, lastChild *html.Node, s *strings.Builder) {
	if lastChild == n {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			s.Write([]byte(c.Data))
		}
		getTextInElement(c, lastChild, s)
	}

}

func hrefVal(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return ""

}

func getAtag(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		s := new(strings.Builder)
		l := new(Link)
		l.Href = hrefVal(n)
		getTextInElement(n, n.LastChild, s)
		tmp := s.String()
		tmp = strings.Join(strings.Fields(tmp), " ")
		l.Text = tmp
		res = append(res, *l)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getAtag(c)
	}
}
