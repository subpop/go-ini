package ini

type invalidKeyErr struct {
	err string
}

func (e *invalidKeyErr) Error() string {
	return "invalid key: " + e.err
}

type parseTree struct {
	global   section
	sections map[string][]section
}

func newParseTree() parseTree {
	return parseTree{
		global:   newSection(""),
		sections: make(map[string][]section),
	}
}

func (p *parseTree) add(s section) {
	sections, ok := p.sections[s.name]
	if !ok {
		sections = make([]section, 0)
	}
	sections = append(sections, s)
	p.sections[s.name] = sections
}

func (p *parseTree) get(name string) ([]section, error) {
	if name == "" {
		return nil, &invalidKeyErr{"section name cannot be empty"}
	}
	if name == "*" {
		sections := make([]section, 0)
		for _, v := range p.sections {
			sections = append(sections, v...)
		}
		return sections, nil
	}
	sections, ok := p.sections[name]
	if !ok {
		sections = make([]section, 0)
		p.sections[name] = sections
	}
	return sections, nil
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

func (s *section) add(p property) {
	s.props[p.key] = p
}

func (s *section) get(key string) (*property, error) {
	if key == "" {
		return nil, &invalidKeyErr{"property key cannot be empty"}
	}
	prop, ok := s.props[key]
	if !ok {
		prop = newProperty(key)
		s.props[key] = prop
	}
	return &prop, nil
}

type property struct {
	key  string
	vals map[string][]string
}

func newProperty(key string) property {
	return property{
		key:  key,
		vals: make(map[string][]string),
	}
}

func (p *property) add(key, value string) {
	vals, ok := p.vals[key]
	if !ok {
		vals = make([]string, 0)
	}
	vals = append(vals, value)
	p.vals[key] = vals
}

func (p *property) get(key string) []string {
	vals, ok := p.vals[key]
	if !ok {
		vals = make([]string, 0)
		p.vals[key] = vals
	}
	return vals
}
