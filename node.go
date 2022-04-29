package html

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Node struct {
	node *html.Node
}

func NewNode(n *html.Node) *Node {
	return node(n)
}

func Parse(reader io.Reader) (*Node, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	return &Node{node}, nil
}

func (n *Node) OuterHTML() string {
	var b bytes.Buffer
	html.Render(&b, n.node)
	return b.String()
}

func (n *Node) InnerHTML() string {
	var b bytes.Buffer

	for node := n.node.FirstChild; node != nil; node = node.NextSibling {
		html.Render(&b, node)
	}

	return b.String()
}

func (n *Node) InnerText() string {
	var text string

	for node := n.node.FirstChild; node != nil; node = node.NextSibling {
		if node.Type == html.TextNode {
			text += node.Data
		} else if node.Type == html.ElementNode {
			text = NewNode(node).InnerText()
		}
	}

	return text
}

func (n *Node) Parents(tag string) *Node {
	selector := NewSelector(tag)

	for node := n.node.Parent; node != nil; node = node.Parent {
		if selector.Matches(node) {
			return NewNode(node)
		}
	}

	return nil
}

func (n *Node) Matches(selector string) bool {
	return NewSelector(selector).Matches(n.node)
}

func (n *Node) GetElementByID(id string) *Node {
	stack := []*html.Node{n.node.FirstChild}

	for len(stack) > 0 {
		x := stack[0]
		stack = stack[1:]

		for node := x; node != nil; node = node.NextSibling {
			if node.Type == html.ElementNode || node.Type == html.DocumentNode {
				for _, attr := range node.Attr {
					if strings.EqualFold(attr.Key, "id") {
						if attr.Val == id {
							return NewNode(node)
						}

						break
					}
				}

				stack = append(stack, node.FirstChild)
			}
		}
	}
	return nil
}

func (n *Node) QuerySelector(selector string, nested ...string) *Node {
	s := NewSelector(selector)
	stack := []*html.Node{n.node.FirstChild}

	for len(stack) > 0 {
		x := stack[0]
		stack = stack[1:]

		for node := x; node != nil; node = node.NextSibling {
			if s.Matches(node) {
				if len(nested) > 0 {
					return NewNode(node).QuerySelector(nested[0], nested[1:]...)
				}

				return NewNode(node)
			}

			if node.Type == html.ElementNode || node.Type == html.DocumentNode {
				stack = append(stack, node.FirstChild)
			}
		}
	}
	return nil
}

func (n *Node) QuerySelectorAll(selector string, nested ...string) []*Node {
	var res []*Node
	s := NewSelector(selector)
	stack := []*html.Node{n.node.FirstChild}

	for len(stack) > 0 {
		x := stack[0]
		stack = stack[1:]

		for node := x; node != nil; node = node.NextSibling {
			if node.Type == html.ElementNode || node.Type == html.DocumentNode {
				if s.Matches(node) {
					if len(nested) > 0 {
						res = append(res, NewNode(node).QuerySelectorAll(nested[0], nested[1:]...)...)
					} else {
						res = append(res, NewNode(node))
					}
				}

				stack = append(stack, node.FirstChild)
			}
		}
	}

	return res
}

func (n *Node) GetElementsByTagName(tag string) []*Node {
	var res []*Node
	stack := []*html.Node{n.node.FirstChild}

	for len(stack) > 0 {
		x := stack[0]
		stack = stack[1:]

		for node := x; node != nil; node = node.NextSibling {
			if node.Type == html.ElementNode || node.Type == html.DocumentNode {
				if node.Data == tag {
					res = append(res, NewNode(node))
				}

				stack = append(stack, node.FirstChild)
			}
		}
	}

	return res
}

func (n *Node) GetElementsByClassName(class string) []*Node {
	var res []*Node
	stack := []*html.Node{n.node.FirstChild}

	for len(stack) > 0 {
		x := stack[0]
		stack = stack[1:]

		for node := x; node != nil; node = node.NextSibling {
			if node.Type == html.ElementNode || node.Type == html.DocumentNode {
				for _, attr := range node.Attr {
					if strings.EqualFold(attr.Key, "class") {
						if attr.Val == class || strings.HasPrefix(attr.Val, class+" ") || strings.HasSuffix(attr.Val, " "+class) || strings.Contains(attr.Val, " "+class+" ") {
							res = append(res, NewNode(node))
						}
						break
					}
				}

				stack = append(stack, node.FirstChild)
			}
		}
	}

	return res
}

func (n *Node) IsElement() bool {
	return n.node.Type == html.ElementNode
}

func (n *Node) IsText() bool {
	return n.node.Type == html.TextNode
}

func (n *Node) Next(selector string) *Node {
	s := NewSelector(selector)

	for node := n.node.NextSibling; node != nil; node = node.NextSibling {
		if s.Matches(node) {
			return NewNode(node)
		}
	}

	return nil
}

func (n *Node) PrevSibling() *Node {
	return node(n.node.PrevSibling)
}

func (n *Node) NextSibling() *Node {
	return node(n.node.NextSibling)
}

func (n *Node) NextElement() *Node {
	for n := n.node.NextSibling; n != nil; n = n.NextSibling {
		if n.Type == html.ElementNode {
			return NewNode(n)
		}
	}

	return nil
}

func (n *Node) PrevElement() *Node {
	for n := n.node.PrevSibling; n != nil; n = n.PrevSibling {
		if n.Type == html.ElementNode {
			return NewNode(n)
		}
	}

	return nil
}

func (n *Node) Previous(selector string) *Node {
	s := NewSelector(selector)

	for node := n.node.PrevSibling; node != nil; node = node.PrevSibling {
		if s.Matches(node) {
			return NewNode(node)
		}
	}

	return nil
}

func (n *Node) Between(from, to *Node) bool {
	if from == nil {
		return false
	}

	for node := from.node.NextSibling; node != nil && (to == nil || node != to.node); node = node.NextSibling {
		if node == n.node {
			return true
		}
	}

	return false
}

func (n *Node) Parent() *Node {
	return node(n.node.Parent)
}

func (n *Node) FirstChild() *Node {
	return node(n.node.FirstChild)
}

func (n *Node) TagName() string {
	if n.node.Type == html.ElementNode {
		return strings.ToLower(n.node.Data)
	}

	return ""
}

func (n *Node) LookupAttribute(name string) (string, bool) {
	for _, a := range n.node.Attr {
		if strings.EqualFold(a.Key, name) {
			return a.Val, true
		}
	}

	return "", false
}

func (n *Node) Attribute(name string) string {
	res, _ := n.LookupAttribute(name)
	return res
}

func (n *Node) String() string {
	return n.OuterHTML()
}

func (n *Node) HTMLNode() *html.Node {
	return n.node
}

func node(node *html.Node) *Node {
	if node == nil {
		return nil
	}

	return &Node{node}
}
