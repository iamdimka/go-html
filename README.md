# html

little util for operating html elements

## Install

`go get -u github.com/iamdimka/go-html`

## API

```go
func NewNode(*html.Node) *Node
func Parse(io.Reader) (*Node, error)

func (*Node) Attribute(name string) bool
func (*Node) Between(from, to *Node) bool
func (*Node) FirstChild() *Node
func (*Node) GetElementByID(id string) *Node
func (*Node) GetElementsByTagName(tag string) []*Node
func (*Node) GetElementsByClassName(tag string) []*Node
func (*Node) HTMLNode() *html.Node
func (*Node) IsElement() string
func (*Node) IsText() string
func (*Node) InnerHTML() string
func (*Node) InnerText() string
func (*Node) Matches(selector string) bool
func (*Node) LookupAttribute(name string) (string, bool)
func (*Node) Next(selector string) *Node
func (*Node) NextElement() *Node
func (*Node) NextSibling() *Node
func (*Node) Parent() *Node
func (*Node) Previous(selector string) *Node
func (*Node) PrevElement() *Node
func (*Node) PrevSibling() *Node
func (*Node) TagName() string
func (*Node) OuterHTML() string
func (*Node) Parents(selector string) *Node
func (*Node) QuerySelector(selector string, children ...string) *Node
func (*Node) QuerySelectorAll(selector string, children ...string) []*Node
```

```go
func NewSelector(query string) *Selector
func (*Selector) Matches(*html.Node) bool
```

**Supported query selector**

`div#id.class1.class2[disabled][data-some=yes]`