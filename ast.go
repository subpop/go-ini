package ini

import (
	"fmt"
)

type missingSectionErr struct {
	name string
}

func (e *missingSectionErr) Error() string {
	return fmt.Sprintf("parse error: missing section name %q", e.name)
}

type missingPropertyErr struct {
	key string
}

func (e *missingPropertyErr) Error() string {
	return fmt.Sprintf("parse error: missing property key %q", e.key)
}

type missingSubkeyErr struct {
	p      property
	subkey string
}

func (e *missingSubkeyErr) Error() string {
	return "property '" + e.p.key + "' missing subkey '" + e.subkey + "'"
}

type parseTree map[string][]section

func newParseTree() parseTree {
	return parseTree{
		"": []section{newSection("")},
	}
}

func (p parseTree) add(s section) {
	sections, ok := p[s.name]
	if !ok {
		sections = make([]section, 0)
	}
	sections = append(sections, s)
	p[s.name] = sections
}

func (p parseTree) get(name string) ([]section, error) {
	sections, ok := p[name]
	if !ok {
		return nil, &missingSectionErr{name}
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

func (s section) add(p property) {
	s.props[p.key] = p
}

func (s section) get(key string) (*property, error) {
	prop, ok := s.props[key]
	if !ok {
		return nil, &missingPropertyErr{key}
	}
	return &prop, nil
}

type property struct {
	key  string
	vals map[string][]string
}

func newProperty(key string) property {
	return property{
		key: key,
		vals: map[string][]string{
			"": []string{},
		},
	}
}

func (p *property) append(subkey string, values ...string) {
	v, ok := p.vals[subkey]
	if !ok {
		v = make([]string, 0)
	}
	v = append(v, values...)
	p.vals[subkey] = v
}

func (p property) values(subkey string) ([]string, error) {
	values, ok := p.vals[subkey]
	if !ok {
		return nil, &missingSubkeyErr{p, subkey}
	}
	return values, nil
}
