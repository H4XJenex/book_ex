package notrelevant

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

const (
	styleTag  = "style"
	scriptTag = "script"
)

func fetch(url string) (*html.Node, error) {

	//if !strings.HasPrefix(url, "https://") {
	//	url = "https://" + url
	//}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	//_, err = io.Copy(os.Stdout, resp.Body)
	//b, err := io.ReadAll(resp.Body)

	b, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}

func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data)
		fmt.Println(stack)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
}

func tagCount(tags map[string]int, n *html.Node) {

	if n.Type == html.ElementNode {
		tags[n.Data]++
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tagCount(tags, c)
	}
}

func tagContent(n *html.Node) {
	if n.Type == html.TextNode && n.Parent.Data != scriptTag && n.Parent.Data != styleTag {

		text := strings.TrimSpace(n.Data)
		if text != "" {
			fmt.Println(text)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tagContent(c)
	}

}

func NewReader(s string) {

}
