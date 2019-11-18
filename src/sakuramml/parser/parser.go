package parser

import (
	"fmt"
	"sakuramml/node"
	"sakuramml/token"
)

// Parser struct
type Parser struct {
	desk token.Desk
	Top  *node.Node
	Last *node.Node
}

// Init parser
func (p *Parser) Init(tokens token.Tokens) {
	p.desk = token.NewDesk(tokens)
	nop := node.NewNop()
	p.Top = nop
	p.Last = p.Top
}

// NewParser func
func NewParser(tokens token.Tokens) *Parser {
	p := Parser{}
	p.Init(tokens)
	return &p
}

func (p *Parser) appendNode(n *node.Node) {
	if n == nil {
		return
	}
	p.Last.Next = n
	p.Last = n
}

func (p *Parser) readNote(t *token.Token) *node.Node {
	n := node.NewNoteOn(t.Label)
	return n
}

func (p *Parser) readValue() *node.Node {
	if p.desk.IsLabel("=") {
		p.desk.Next()
	}
	ct := p.desk.Peek()
	if p.desk.IsType(token.Number) {
		nn := node.NewNumber(ct.Label)
		p.desk.Next()
		return nn
	}
	panic("not implement")
}

func (p *Parser) readTrack() *node.Node {
	no := p.readValue()
	return node.NewTrack(no)
}

func (p *Parser) readWord() *node.Node {
	t := p.desk.Next()
	switch t.Label {
	case "c":
		return p.readNote(t)
	case "d":
		return p.readNote(t)
	case "e":
		return p.readNote(t)
	case "f":
		return p.readNote(t)
	case "g":
		return p.readNote(t)
	case "a":
		return p.readNote(t)
	case "b":
		return p.readNote(t)
	case "TR":
		return p.readTrack()
	}
	return nil
}

// Parse func
func (p *Parser) Parse() (*node.Node, error) {
	for p.desk.HasNext() {
		t := p.desk.Peek()
		fmt.Printf("Parse %s\n", t.Label)
		if t.Type == token.Word {
			p.appendNode(p.readWord())
			continue
		}
		e := fmt.Errorf("[ERROR] (%d) not implements : %s ", p.desk.Peek().Line, p.desk.Peek().Label)
		return p.Top, e
	}
	return p.Top, nil
}

// Parse convert to AST
func Parse(tokens token.Tokens) (*node.Node, error) {
	return NewParser(tokens).Parse()
}
