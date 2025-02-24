package html

import (
	"fmt"
	"strings"
	"unicode"
)

// Attribute represents an HTML attribute with key and value.
type Attribute struct {
	Key   string
	Value string
}

// Node is the interface for nodes in our AST.
type Node interface {
	Dump() string
}

type NodeList []Node

func (nodeList NodeList) Dump() string {
	var builder strings.Builder
	for _, node := range nodeList {
		builder.WriteString(node.Dump())
	}
	return builder.String()
}

// RawNode holds unchanged text.
type RawNode struct {
	Text string
	Line int // added field
}

// Dump returns the raw text.
func (r *RawNode) Dump() string {
	return r.Text
}

// ElementNode represents an HTML element.
type ElementNode struct {
	Tag         string
	Attributes  []Attribute
	Children    NodeList
	SelfClosing bool
	Line        int // added field
}

// Dump returns the HTML representation of the element and its children.
func (e *ElementNode) Dump() string {
	var builder strings.Builder
	builder.WriteString("<" + e.Tag)
	// Use attributes slice to preserve order.
	for _, attr := range e.Attributes {
		if attr.Value == "" {
			builder.WriteString(" " + attr.Key)
		} else {
			builder.WriteString(" " + attr.Key + "=\"" + attr.Value + "\"")
		}
	}
	if e.SelfClosing {
		builder.WriteString("/>")
	} else {
		builder.WriteString(">")
		for _, child := range e.Children {
			builder.WriteString(child.Dump())
		}
		builder.WriteString("</" + e.Tag + ">")
	}
	return builder.String()
}

// Parser holds the state for our simple parser.
type Parser struct {
	input  string
	pos    int
	length int
}

// NewParser creates a new parser for the given input.
func NewParser(input string) (NodeList, error) {
	p := &Parser{input: input, pos: 0, length: len(input)}

	return p.parseNodes("")
}

// current returns the current byte (or zero if at the end).
func (p *Parser) current() byte {
	if p.pos >= p.length {
		return 0
	}
	return p.input[p.pos]
}

// peek returns the next n characters (or what remains).
func (p *Parser) peek(n int) string {
	if p.pos+n > p.length {
		return p.input[p.pos:]
	}
	return p.input[p.pos : p.pos+n]
}

// skipWhitespace advances the position over any whitespace.
func (p *Parser) skipWhitespace() {
	for p.pos < p.length &&
		(p.input[p.pos] == ' ' || p.input[p.pos] == '\n' ||
			p.input[p.pos] == '\r' || p.input[p.pos] == '\t') {
		p.pos++
	}
}

// Helper to get line number at a given position.
func (p *Parser) getLineAt(pos int) int {
	return strings.Count(p.input[:pos], "\n") + 1
}

// skipComment skips an HTML comment starting at "<!--" until "-->".
func (p *Parser) skipComment() error {
	if p.peek(4) != "<!--" {
		return nil
	}
	// Skip the opening "<!--"
	p.pos += 4
	// Find the closing "-->"
	idx := strings.Index(p.input[p.pos:], "-->")
	if idx == -1 {
		return fmt.Errorf("unterminated comment starting at pos %d", p.pos-4)
	}
	p.pos += idx + 3 // skip past "-->"
	return nil
}

// parseNodes parses a list of nodes until an optional stop tag (used for element children).
func (p *Parser) parseNodes(stopTag string) (NodeList, error) {
	var nodes NodeList
	rawStart := p.pos

	for p.pos < p.length {
		// Check for comment start and skip it.
		if p.peek(4) == "<!--" {
			if p.pos > rawStart {
				text := p.input[rawStart:p.pos]
				if text != "" {
					nodes = append(nodes, &RawNode{
						Text: text,
						Line: p.getLineAt(rawStart), // added line attribute
					})
				}
			}
			if err := p.skipComment(); err != nil {
				return nodes, err
			}
			rawStart = p.pos
			continue
		}
		// If we’re about to hit a closing tag for the current element, break.
		if p.current() == '<' && p.peek(2) == "</" {
			// Save position to check tag name.
			savedPos := p.pos
			p.pos += 2
			p.skipWhitespace()
			closingTag := p.parseTagName()
			// Reset position so the caller can see the closing tag.
			p.pos = savedPos
			if stopTag != "" && closingTag == stopTag {
				break
			}
		}

		// If we see a '<', then try to parse an element node.
		if p.current() == '<' {
			// If any raw text is accumulated, add it as a RawNode.
			if p.pos > rawStart {
				text := p.input[rawStart:p.pos]
				if text != "" {
					nodes = append(nodes, &RawNode{
						Text: text,
						Line: p.getLineAt(rawStart), // added line attribute
					})
				}
			}
			element, err := p.parseElement()
			if err != nil {
				return nodes, err
			}
			nodes = append(nodes, element)
			rawStart = p.pos // mark new raw text start
		} else {
			p.pos++
		}
	}
	// Append any remaining raw text.
	if rawStart < p.pos {
		text := p.input[rawStart:p.pos]
		if text != "" {
			nodes = append(nodes, &RawNode{
				Text: text,
				Line: p.getLineAt(rawStart), // added line attribute
			})
		}
	}
	return nodes, nil
}

// isVoidElement returns true if the tag is a void element (e.g., <br> does not require a closing tag)
func isVoidElement(tag string) bool {
	switch strings.ToLower(tag) {
	case "area", "base", "br", "col", "embed", "hr", "img", "input", "keygen", "link", "meta", "param", "source", "track", "wbr":
		return true
	}
	return false
}

// parseElement parses an HTML element starting at the current position (assumes a '<').
func (p *Parser) parseElement() (Node, error) {
	// Record start position for line number.
	startPos := p.pos
	if p.current() != '<' {
		return nil, fmt.Errorf("expected '<' at pos %d", p.pos)
	}
	p.pos++ // skip '<'
	p.skipWhitespace()

	tagName := p.parseTagName()
	if tagName == "" {
		return nil, fmt.Errorf("empty tag name at pos %d", p.pos)
	}

	node := &ElementNode{
		Tag:        tagName,
		Attributes: []Attribute{},
		Children:   NodeList{},
		Line:       p.getLineAt(startPos), // assign starting line
	}

	// Parse element attributes.
	for p.pos < p.length {
		p.skipWhitespace()
		if p.current() == '>' || (p.current() == '/' && p.peek(2) == "/>") {
			break
		}
		attrName := p.parseAttrName()
		if attrName == "" {
			break
		}
		p.skipWhitespace()
		var attrVal string
		if p.current() == '=' {
			p.pos++ // skip '='
			p.skipWhitespace()
			attrVal = p.parseAttrValue()
		}
		// Append attribute preserving order.
		node.Attributes = append(node.Attributes, Attribute{Key: attrName, Value: attrVal})
	}

	// Check for self-closing tag.
	if p.current() == '/' {
		p.pos++ // skip '/'
		if p.current() != '>' {
			return nil, fmt.Errorf("expected '>' after '/' at pos %d", p.pos)
		}
		p.pos++ // skip '>'
		node.SelfClosing = true
		return node, nil
	}
	if p.current() == '>' {
		p.pos++ // skip '>'
		if isVoidElement(tagName) {
			node.SelfClosing = true
			return node, nil
		}
	} else {
		return nil, fmt.Errorf("expected '>' at pos %d", p.pos)
	}

	// Parse children until the corresponding closing tag.
	children, err := p.parseElementChildren(node.Tag)
	if err != nil {
		return nil, err
	}
	node.Children = children

	return node, nil
}

// parseElementChildren parses the child nodes of an element until the closing tag is reached.
func (p *Parser) parseElementChildren(tag string) (NodeList, error) {
	var children NodeList
	rawStart := p.pos

	for p.pos < p.length {
		// Check for comment and skip.
		if p.peek(4) == "<!--" {
			if p.pos > rawStart {
				text := p.input[rawStart:p.pos]
				if text != "" {
					children = append(children, &RawNode{
						Text: text,
						Line: p.getLineAt(rawStart), // added line attribute
					})
				}
			}
			if err := p.skipComment(); err != nil {
				return children, err
			}
			rawStart = p.pos
			continue
		}
		// Check for a closing tag.
		if p.current() == '<' && p.peek(2) == "</" {
			savedPos := p.pos
			p.pos += 2 // skip "</"
			p.skipWhitespace()
			closingTag := p.parseTagName()
			p.skipWhitespace()
			if p.current() == '>' {
				p.pos++ // skip '>'
			} else {
				return children,
					fmt.Errorf("expected '>' for closing tag at pos %d", p.pos)
			}
			if closingTag == tag {
				// Add any raw text before the closing tag.
				if rawStart < savedPos {
					text := p.input[rawStart:savedPos]
					if text != "" {
						children = append(children, &RawNode{
							Text: text,
							Line: p.getLineAt(rawStart), // added line attribute
						})
					}
				}
				return children, nil
			} else {
				// Not the matching closing tag; reset and continue.
				p.pos = savedPos
			}
		}

		if p.current() == '<' {
			if p.pos > rawStart {
				text := p.input[rawStart:p.pos]
				if text != "" {
					children = append(children, &RawNode{
						Text: text,
						Line: p.getLineAt(rawStart), // added line attribute
					})
				}
			}
			child, err := p.parseElement()
			if err != nil {
				return children, err
			}
			children = append(children, child)
			rawStart = p.pos
		} else {
			p.pos++
		}
	}
	return children, nil
}

// parseTagName parses a tag or attribute name (letters, digits, '-' and ':').
func (p *Parser) parseTagName() string {
	start := p.pos
	for p.pos < p.length {
		c := p.input[p.pos]
		if unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || c == '-' || c == ':' {
			p.pos++
		} else {
			break
		}
	}
	return p.input[start:p.pos]
}

// parseAttrName parses an attribute name.
func (p *Parser) parseAttrName() string {
	start := p.pos
	// Accept characters until whitespace, '=', '>', or '/'.
	for p.pos < p.length {
		c := p.input[p.pos]
		if c == ' ' || c == '\n' || c == '\r' || c == '\t' ||
			c == '=' || c == '>' || c == '/' {
			break
		}
		p.pos++
	}
	return p.input[start:p.pos]
}

// parseAttrValue parses an attribute value (expects a quoted string).
func (p *Parser) parseAttrValue() string {
	if p.current() == '"' {
		p.pos++ // skip opening "
		start := p.pos
		for p.pos < p.length && p.current() != '"' {
			p.pos++
		}
		val := p.input[start:p.pos]
		if p.current() == '"' {
			p.pos++ // skip closing "
		}
		return val
	}
	// Allow unquoted values.
	start := p.pos
	for p.pos < p.length &&
		p.current() != ' ' && p.current() != '>' {
		p.pos++
	}
	return p.input[start:p.pos]
}

func TraverseNode(n NodeList, f func(*ElementNode)) {
	for _, node := range n {
		switch node := node.(type) {
		case *ElementNode:
			f(node)
			for _, child := range node.Children {
				TraverseNode(NodeList{child}, f)
			}
		}
	}
}
