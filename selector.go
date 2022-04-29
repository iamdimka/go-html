package html

import (
	"bytes"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

type Selector struct {
	tag     string
	classes []string
	attr    map[string]string
}

func NewSelector(query string) *Selector {
	s := &Selector{
		attr: make(map[string]string),
	}

	for {
		idx := strings.LastIndexAny(query, "#.[")
		if idx < 0 {
			break
		}

		switch query[idx] {
		case '#':
			s.attr["id"] = query[idx+1:]

		case '.':
			s.classes = append(s.classes, query[idx+1:])

		case '[':
			end := strings.LastIndexByte(query, ']')
			if end < 0 {
				panic("end of tag ] not found")
			}

			data := query[idx+1 : end]
			eq := strings.IndexByte(data, '=')
			if eq < 0 {
				s.attr[strings.ToLower(data)] = "[]"
			} else {
				s.attr[strings.ToLower(data[:eq])] = data[eq+1:]
			}
		}

		query = query[:idx]
	}

	s.tag = strings.ToLower(query)
	return s
}

func (s *Selector) Matches(node *html.Node) bool {
	if node.Type != html.ElementNode {
		return false
	}

	if s.tag != "" && strings.ToLower(node.Data) != s.tag {
		return false
	}

	if len(s.classes) > 0 {
		found := false
		for _, a := range node.Attr {
			if !strings.EqualFold(a.Key, "class") {
				continue
			}

			found = true
			cl := sort.StringSlice(strings.Split(a.Val, " "))
			for _, c := range s.classes {
				if sort.SearchStrings(cl, c) == -1 {
					return false
				}
			}
			break
		}

		if !found {
			return false
		}
	}

	if len(s.attr) > 0 {
		for _, a := range node.Attr {
			if val, ok := s.attr[strings.ToLower(a.Key)]; ok {
				if val != "[]" && val != a.Val {
					return false
				}
			}
		}
	}

	return true
}

func (s *Selector) String() string {
	var b bytes.Buffer
	if s.tag != "" {
		b.WriteString(s.tag)
	}

	id := s.attr["id"]
	if id != "" {
		b.WriteByte('#')
		b.WriteString(id)
	}

	for _, c := range s.classes {
		b.WriteByte('.')
		b.WriteString(c)
	}

	for key, val := range s.attr {
		if key == "id" {
			continue
		}

		b.WriteByte('[')
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteString(val)
		b.WriteByte(']')
	}

	return b.String()
}
