package ini

import (
	"errors"
	"fmt"
)

type unexpectedTokenErr struct {
	got  token
	want token
}

func (e *unexpectedTokenErr) Error() string {
	return fmt.Sprintf("unexpected token: %v, want %v", e.got, e.want)
}

type invalidPropertyErr struct {
	p property
}

func (e *invalidPropertyErr) Error() string {
	return "invalid property: " + e.p.key
}

type parser struct {
	tree parseTree
	l    *lexer
	tok  token
	prev *token
}

func newParser(data []byte) *parser {
	p := parser{
		tree: newParseTree(),
		l:    lex(string(data)),
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
			return errors.New(p.tok.val)
		case tokenSection:
			sec := newSection(p.tok.val)
			if err := p.parseSection(&sec); err != nil {
				return err
			}
			p.tree.add(sec)
		case tokenPropKey:
			prop, err := p.tree.global.get(p.tok.val)
			if err != nil {
				return err
			}
			if err := p.parseProperty(prop); err != nil {
				return err
			}
			p.tree.global.add(*prop)
		case tokenComment:
			continue
		default:
			return &unexpectedTokenErr{got: p.tok}
		}
	}
}

func (p *parser) parseSection(out *section) error {
	name := p.tok.val
	out.name = name

	for {
		p.nextToken()
		switch p.tok.typ {
		case tokenError:
			return errors.New(p.tok.val)
		case tokenPropKey:
			prop, err := out.get(p.tok.val)
			if err != nil {
				return err
			}
			if err := p.parseProperty(prop); err != nil {
				return err
			}
			out.add(*prop)
		case tokenSection:
			// we've parsed too far; backup so we can parse the next section
			p.backup()
			return nil
		case tokenComment:
			continue
		default:
			return nil
		}
	}
}

func (p *parser) parseProperty(out *property) error {
	key := p.tok.val
	subkey := ""

	p.nextToken()
	if p.tok.typ == tokenMapKey {
		subkey = p.tok.val
		p.nextToken()
	}

	if p.tok.typ != tokenAssignment {
		return &unexpectedTokenErr{
			got: p.tok,
			want: token{
				typ: tokenAssignment,
				val: "=",
			},
		}
	}

	p.nextToken()
	if p.tok.typ != tokenPropValue {
		return &unexpectedTokenErr{
			got: p.tok,
			want: token{
				typ: tokenPropValue,
			},
		}
	}
	val := p.tok.val

	out.key = key
	out.add(subkey, val)

	return nil
}
