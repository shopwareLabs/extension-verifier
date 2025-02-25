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

// CommentNode represents an HTML comment
type CommentNode struct {
	Text string
	Line int
}

// Dump returns the comment text with HTML comment syntax
func (c *CommentNode) Dump() string {
	return "<!-- " + c.Text + " -->"
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
	return e.dump(0)
}

// dump formats the element with the given indentation level
func (e *ElementNode) dump(indent int) string {
	var builder strings.Builder

	// Add initial indentation
	for i := 0; i < indent; i++ {
		builder.WriteString("\t")
	}

	builder.WriteString("<" + e.Tag)

	// Add attributes on new lines with indentation if there are multiple attributes
	if len(e.Attributes) > 0 {
		if len(e.Attributes) == 1 {
			// Single attribute on same line
			attr := e.Attributes[0]
			builder.WriteString(" ")
			if attr.Value == "" {
				builder.WriteString(attr.Key)
			} else {
				builder.WriteString(attr.Key + "=\"" + attr.Value + "\"")
			}
		} else {
			// Multiple attributes on new lines
			for _, attr := range e.Attributes {
				builder.WriteString("\n")
				for i := 0; i < indent+1; i++ {
					builder.WriteString("\t")
				}
				if attr.Value == "" {
					builder.WriteString(attr.Key)
				} else {
					builder.WriteString(attr.Key + "=\"" + attr.Value + "\"")
				}
			}
		}
	}

	if e.SelfClosing && len(e.Children) == 0 {
		if len(e.Attributes) > 1 {
			builder.WriteString("\n")
			for i := 0; i < indent; i++ {
				builder.WriteString("\t")
			}
		}
		builder.WriteString("/>")
		return builder.String()
	}

	// Close opening tag
	if len(e.Attributes) > 1 {
		builder.WriteString("\n")
		for i := 0; i < indent; i++ {
			builder.WriteString("\t")
		}
	}
	builder.WriteString(">")

	// Special case: if there's only one child and it's a text node or comment, don't add newlines
	if len(e.Children) == 1 {
		if _, ok := e.Children[0].(*RawNode); ok {
			builder.WriteString(e.Children[0].Dump())
			builder.WriteString("</" + e.Tag + ">")
			return builder.String()
		}
		if _, ok := e.Children[0].(*CommentNode); ok {
			builder.WriteString(e.Children[0].Dump())
			builder.WriteString("</" + e.Tag + ">")
			return builder.String()
		}
	}

	// Special case: if all children are comments or text nodes, keep them on same line
	allSimpleNodes := true
	for _, child := range e.Children {
		if _, ok := child.(*RawNode); !ok {
			if _, ok := child.(*CommentNode); !ok {
				allSimpleNodes = false
				break
			}
		}
	}

	if allSimpleNodes && len(e.Children) > 0 {
		for _, child := range e.Children {
			builder.WriteString(child.Dump())
		}
		builder.WriteString("</" + e.Tag + ">")
		return builder.String()
	}

	// Add children with increased indentation
	if len(e.Children) > 0 {
		builder.WriteString("\n")
		for i, child := range e.Children {
			if elementChild, ok := child.(*ElementNode); ok {
				builder.WriteString(elementChild.dump(indent + 1))
			} else {
				// For text nodes and comments, indent and write directly
				for i := 0; i < indent+1; i++ {
					builder.WriteString("\t")
				}
				builder.WriteString(child.Dump())
			}
			if i < len(e.Children)-1 {
				builder.WriteString("\n")
			}
		}
		builder.WriteString("\n")
		for i := 0; i < indent; i++ {
			builder.WriteString("\t")
		}
	}
	builder.WriteString("</" + e.Tag + ">")
	return builder.String()
}

// TwigBlockNode represents a twig block
type TwigBlockNode struct {
	Name     string
	Children NodeList
	Line     int
}

// Dump returns the twig block with proper formatting
func (t *TwigBlockNode) Dump() string {
	var builder strings.Builder
	builder.WriteString("{% block " + t.Name + " %}")
	if len(t.Children) > 0 {
		builder.WriteString("\n")
		for _, child := range t.Children {
			if elementChild, ok := child.(*ElementNode); ok {
				builder.WriteString(elementChild.dump(1))
			} else {
				builder.WriteString("\t")
				builder.WriteString(child.Dump())
			}
			builder.WriteString("\n")
		}
	}
	builder.WriteString("{% endblock %}")
	return builder.String()
}

// ParentNode represents a twig parent() call
type ParentNode struct {
	Line int
}

func (p *ParentNode) Dump() string {
	return "{% parent() %}"
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

// parseComment parses an HTML comment and returns a CommentNode
func (p *Parser) parseComment() (*CommentNode, error) {
	if p.peek(4) != "<!--" {
		return nil, fmt.Errorf("expected comment at pos %d", p.pos)
	}
	startPos := p.pos
	p.pos += 4 // skip "<!--"

	start := p.pos
	idx := strings.Index(p.input[p.pos:], "-->")
	if idx == -1 {
		return nil, fmt.Errorf("unterminated comment starting at pos %d", startPos)
	}

	commentText := strings.TrimSpace(p.input[start : start+idx])
	p.pos += idx + 3 // skip past "-->"

	return &CommentNode{
		Text: commentText,
		Line: p.getLineAt(startPos),
	}, nil
}

// parseNodes parses a list of nodes until an optional stop tag (used for element children).
func (p *Parser) parseNodes(stopTag string) (NodeList, error) {
	var nodes NodeList
	rawStart := p.pos

	for p.pos < p.length {
		// Check for endblock if we're parsing twig block children
		if stopTag == "" && p.peek(2) == "{%" {
			peek := p.input[p.pos:]
			if strings.HasPrefix(peek, "{% endblock") {
				break
			}
		}

		if p.peek(2) == "{%" {
			if p.pos > rawStart {
				text := p.input[rawStart:p.pos]
				if text != "" {
					nodes = append(nodes, &RawNode{
						Text: text,
						Line: p.getLineAt(rawStart),
					})
				}
			}

			// Try parsing twig directives first
			directive, err := p.parseTwigDirective()
			if err != nil {
				return nodes, err
			}
			if directive != nil {
				nodes = append(nodes, directive)
				rawStart = p.pos
				continue
			}

			// If not a directive, try parsing as a block
			startPos := p.pos
			block, err := p.parseTwigBlock()
			if err != nil {
				return nodes, err
			}
			if block != nil {
				nodes = append(nodes, block)
				rawStart = p.pos
			} else {
				// If it wasn't a block, reset position and continue as raw text
				p.pos = startPos
			}
			continue
		}

		if p.peek(4) == "<!--" {
			if p.pos > rawStart {
				text := p.input[rawStart:p.pos]
				if text != "" {
					nodes = append(nodes, &RawNode{
						Text: text,
						Line: p.getLineAt(rawStart),
					})
				}
			}
			comment, err := p.parseComment()
			if err != nil {
				return nodes, err
			}
			nodes = append(nodes, comment)
			rawStart = p.pos
			continue
		}

		// If we're about to hit a closing tag for the current element, break.
		if p.current() == '<' && p.peek(2) == "</" {
			savedPos := p.pos
			p.pos += 2
			p.skipWhitespace()
			closingTag := p.parseTagName()
			p.pos = savedPos
			if stopTag != "" && closingTag == stopTag {
				break
			}
		}

		if p.current() == '<' && p.peek(2) != "<!--" {
			if p.pos > rawStart {
				text := p.input[rawStart:p.pos]
				if text != "" {
					nodes = append(nodes, &RawNode{
						Text: text,
						Line: p.getLineAt(rawStart),
					})
				}
			}
			element, err := p.parseElement()
			if err != nil {
				return nodes, err
			}
			nodes = append(nodes, element)
			rawStart = p.pos
		} else {
			p.pos++
		}
	}

	if rawStart < p.pos {
		text := p.input[rawStart:p.pos]
		if text != "" {
			nodes = append(nodes, &RawNode{
				Text: text,
				Line: p.getLineAt(rawStart),
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
		if p.peek(4) == "<!--" {
			if p.pos > rawStart {
				text := p.input[rawStart:p.pos]
				if text != "" {
					children = append(children, &RawNode{
						Text: text,
						Line: p.getLineAt(rawStart),
					})
				}
			}
			comment, err := p.parseComment()
			if err != nil {
				return children, err
			}
			children = append(children, comment)
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
							Line: p.getLineAt(rawStart),
						})
					}
				}
				return children, nil
			} else {
				// Not the matching closing tag; reset and continue.
				p.pos = savedPos
			}
		}

		if p.current() == '<' && p.peek(2) != "<!--" {
			if p.pos > rawStart {
				text := p.input[rawStart:p.pos]
				if text != "" {
					children = append(children, &RawNode{
						Text: text,
						Line: p.getLineAt(rawStart),
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
	// Accept characters until whitespace, '=', '>', or '/'
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

func (p *Parser) parseTwigDirective() (Node, error) {
	if p.peek(2) != "{%" {
		return nil, nil
	}

	startPos := p.pos
	p.pos += 2 // skip "{%"
	p.skipWhitespace()

	// Check if it's a parent() call
	if strings.HasPrefix(p.input[p.pos:], "parent()") {
		p.pos += 8 // skip "parent()"
		p.skipWhitespace()
		if p.peek(2) != "%}" {
			return nil, fmt.Errorf("unclosed parent directive at pos %d", startPos)
		}
		p.pos += 2 // skip "%}"
		return &ParentNode{Line: p.getLineAt(startPos)}, nil
	}

	// Reset position if it's not a recognized directive
	p.pos = startPos
	return nil, nil
}

func (p *Parser) parseTwigBlock() (Node, error) {
	if p.peek(2) != "{%" {
		return nil, nil
	}

	startPos := p.pos
	p.pos += 2 // skip "{%"
	p.skipWhitespace()

	// Check if it's a block
	if !strings.HasPrefix(p.input[p.pos:], "block") {
		p.pos = startPos
		return nil, nil
	}
	p.pos += 5 // skip "block"
	p.skipWhitespace()

	// Parse block name
	start := p.pos
	for p.pos < p.length && p.current() != '%' && p.current() != ' ' {
		p.pos++
	}
	name := strings.TrimSpace(p.input[start:p.pos])

	// Skip to end of opening tag
	for p.pos < p.length && p.peek(2) != "%}" {
		p.pos++
	}
	if p.peek(2) != "%}" {
		return nil, fmt.Errorf("unclosed block tag at pos %d", startPos)
	}
	p.pos += 2 // skip "%}"

	// Parse children until endblock
	children, err := p.parseNodes("")
	if err != nil {
		return nil, err
	}

	// Look for endblock
	p.skipWhitespace()
	if !strings.HasPrefix(p.input[p.pos:], "{%") {
		return nil, fmt.Errorf("missing endblock at pos %d", p.pos)
	}
	p.pos += 2 // skip "{%"
	p.skipWhitespace()
	
	if !strings.HasPrefix(p.input[p.pos:], "endblock") {
		return nil, fmt.Errorf("missing endblock at pos %d", p.pos)
	}
	p.pos += 8 // skip "endblock"

	// Skip to end of closing tag
	for p.pos < p.length && p.peek(2) != "%}" {
		p.pos++
	}
	if p.peek(2) != "%}" {
		return nil, fmt.Errorf("unclosed endblock tag at pos %d", p.pos)
	}
	p.pos += 2 // skip "%}"

	return &TwigBlockNode{
		Name:     name,
		Children: children,
		Line:     p.getLineAt(startPos),
	}, nil
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
