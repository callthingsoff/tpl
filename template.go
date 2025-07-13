package tpl

import (
	"iter"
	"os"

	"gopkg.in/yaml.v3"
)

// -------- interfaces

// Tpler represents the root.
type Tpler interface {
	GetItems() iter.Seq[Grouper]
}

// Grouper represents the groups.
type Grouper interface {
	GetURL() string
	GetGroup() iter.Seq[JSONPather]
}

// JSONPather represents the basic object.
type JSONPather interface {
	GetID() string
	GetName() string
	GetJSONPath() string
}

// -------- concrete types

// Template holds resources in yaml.
type Template struct {
	Items []itemType `yaml:"items"`
}

func (t Template) GetItems() iter.Seq[Grouper] {
	return func(yield func(Grouper) bool) {
		for _, v := range t.Items {
			if !yield(v) {
				return
			}
		}
	}
}

type groupType struct {
	ID       string `yaml:"id"`
	Name     string `yaml:"name"`
	JSONPath string `yaml:"jsonpath"`
}

func (g groupType) GetID() string {
	return g.ID
}

func (g groupType) GetName() string {
	return g.Name
}

func (g groupType) GetJSONPath() string {
	return g.JSONPath
}

type itemType struct {
	URL   string      `yaml:"url"`
	Group []groupType `yaml:"group"`
}

func (it itemType) GetURL() string {
	return it.URL
}

func (it itemType) GetGroup() iter.Seq[JSONPather] {
	return func(yield func(JSONPather) bool) {
		for _, v := range it.Group {
			if !yield(v) {
				return
			}
		}
	}
}

// -------- parsers

// ParseBytes parses bytes to Template.
func ParseBytes(b []byte) (Tpler, error) {
	var out Template
	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ParseFile parses file to Template.
func ParseFile(filename string) (Tpler, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseBytes(b)
}
