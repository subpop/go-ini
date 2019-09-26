package ini

import (
	"fmt"
	"regexp"
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

func (a ast) getSection(name string) ([]section, error) {
	sections, ok := a[name]
	if !ok {
		return nil, &missingSectionErr{name}
	}
	return sections, nil
}

func (a ast) getSectionMatch(r *regexp.Regexp) []section {
	sections := make([]section, 0)
	for name, section := range a {
		if name == "" {
			continue
		}
		if r.MatchString(name) {
			sections = append(sections, section...)
		}
	}
	return sections
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

func (s section) getProperty(key string) (*property, error) {
	prop, ok := s.props[key]
	if !ok {
		return nil, &missingPropertyErr{key}
	}
	return &prop, nil
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
