package ini

type ast map[string][]section

// newAST returns an AST containing a single-element, global property section.
func newAST() ast {
	return ast{
		"": []section{newSection("")},
	}
}

func (a ast) addSection(s section) {
	sec, ok := a[s.name]
	if !ok {
		sec = make([]section, 0)
	}
	sec = append(sec, s)
	a[s.name] = sec
}

type section struct {
	name  string
	props map[string]property
}

func newSection(name string) section {
	return section{
		name:  name,
		props: make(map[string]property),
	}
}

func (s section) addProperty(p property) {
	prop, ok := s.props[p.key]
	if !ok {
		prop = newProperty(p.key)
	}
	prop.val = append(prop.val, p.val...)

	s.props[p.key] = prop
}

type property struct {
	key string
	val []string
}

func newProperty(key string) property {
	return property{
		key: key,
		val: make([]string, 0),
	}
}

func (p *property) appendVal(s string) {
	p.val = append(p.val, s)
}
