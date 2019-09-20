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
	ast map[string]section
	l   *lexer
	tok token
}

func newParser(data []byte) *parser {
	p := parser{
		ast: map[string]section{
			"": section{
				name:  "",
				props: map[string]property{},
			},
		},
		l: lex(string(data)),
	}
	return &p
}

func (p *parser) nextToken() {
	p.tok = p.l.nextToken()
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
			var sec section
			sec, ok := p.ast[p.tok.val]
			if !ok {
				sec = section{
					name:  p.tok.val,
					props: map[string]property{},
				}
			}
			if err := p.parseSection(&sec); err != nil {
				return err
			}
			p.ast[sec.name] = sec
		case tokenKey:
			var prop property
			prop, ok := p.ast[""].props[p.tok.val]
			if !ok {
				prop = property{
					key: p.tok.val,
					val: []string{},
				}
			}
			if err := p.parseProperty(&prop); err != nil {
				return err
			}
			p.ast[""].props[prop.key] = prop
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
			var prop property
			prop, ok := out.props[p.tok.val]
			if !ok {
				prop = property{
					key: p.tok.val,
					val: []string{},
				}
			}
			if err := p.parseProperty(&prop); err != nil {
				return err
			}
			out.props[prop.key] = prop
		case tokenSection:
			var sec section
			sec, ok := p.ast[p.tok.val]
			if !ok {
				sec = section{
					name:  p.tok.val,
					props: map[string]property{},
				}
			}
			if err := p.parseSection(&sec); err != nil {
				return err
			}
			p.ast[sec.name] = sec
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
