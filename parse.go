package ini

// unexpectedTokenErr describes a token that was not expected by the parser in
// the lexer's current state.
type unexpectedTokenErr struct {
	got token
}

func (e unexpectedTokenErr) Error() string {
	return "unexpected token: " + e.got.val
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

// parse advances the token scanner repeatedly, constructing a parseTree on
// each step through the token stream until an EOF token is encountered.
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
			return &unexpectedTokenErr{p.tok}
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

// parseSection repeatedly advances the token scanner, constructing a section
// parseTree element from the scanned values.
func (p *parser) parseSection(out *section) error {
	name := p.tok.val
	out.name = name

	for {
		p.nextToken()
		switch p.tok.typ {
		case tokenError:
			return &unexpectedTokenErr{got: p.tok}
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

// parseProperty repeatedly advances the token scanner, constructing a property
// parseTree element from the scanned values.
func (p *parser) parseProperty(out *property) error {
	key := p.tok.val
	subkey := ""

	p.nextToken()
	if p.tok.typ == tokenMapKey {
		subkey = p.tok.val
		p.nextToken()
	}

	p.nextToken()
	if p.tok.typ != tokenPropValue {
		return &unexpectedTokenErr{
			got: p.tok,
		}
	}
	val := p.tok.val

	out.key = key
	out.add(subkey, val)

	return nil
}
