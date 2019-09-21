package ini

import (
	"fmt"
)

type errParse struct {
	line int
	col  int
	s    string
}

func (e *errParse) Error() string {
	return fmt.Sprintf("parse:%v:%v: %v", e.line, e.col, e.s)
}

type parser struct {
	ast  ast
	l    *lexer
	tok  token
	prev *token
}

func newParser(data []byte) *parser {
	p := parser{
		ast: newAST(),
		l:   lex(string(data)),
	}
	return &p
}

func (p *parser) nextToken() {
	if p.prev != nil {
		p.tok = *p.prev
		p.prev = nil
	} else {
		p.tok = p.l.nextToken()
	}
}

func (p *parser) backup() {
	p.prev = &p.tok
}

func (p *parser) parse() error {
	for {
		if p.tok.typ == tokenEOF {
			return nil
		}
		p.nextToken()
		switch p.tok.typ {
		case tokenEOF:
			return nil
		case tokenError:
			return &errParse{p.l.line, p.l.col, p.tok.val}
		case tokenSection:
			sec := newSection(p.tok.val)
			if err := p.parseSection(&sec); err != nil {
				return err
			}
			p.ast.addSection(sec)
			p.backup()
		case tokenKey:
			prop := newProperty(p.tok.val)
			if err := p.parseProperty(&prop); err != nil {
				return err
			}
			p.ast[""][0].addProperty(prop)
		}
	}
}

func (p *parser) parseSection(out *section) error {
	name := p.tok.val

	if name != out.name {
		panic(fmt.Sprintf("section name mismatch: expected '%v', got '%v'", name, out.name))
	}

	for {
		if p.tok.typ == tokenEOF {
			return nil
		}
		p.nextToken()
		switch p.tok.typ {
		case tokenEOF:
			return nil
		case tokenError:
			return &errParse{p.l.line, p.l.col, p.tok.val}
		case tokenKey:
			prop := newProperty(p.tok.val)
			if err := p.parseProperty(&prop); err != nil {
				return err
			}
			out.addProperty(prop)
		default:
			return nil
		}
	}
}

func (p *parser) parseProperty(out *property) error {
	key := p.tok.val

	p.nextToken()
	if p.tok.typ != tokenAssignment {
		return &errParse{p.l.line, p.l.col, p.tok.val}
	}

	p.nextToken()
	if p.tok.typ != tokenText {
		return &errParse{p.l.line, p.l.col, p.tok.val}
	}
	val := p.tok.val

	if key != out.key {
		panic(fmt.Sprintf("property key mismatch: expected '%v', got '%v'", key, out.key))
	}

	out.val = append(out.val, val)

	return nil
}
