package tpl

import (
	"os"

	"gopkg.in/yaml.v3"
)

type groupType struct {
	ID       string `yaml:"id"`
	Name     string `yaml:"name"`
	JSONPath string `yaml:"jsonpath"`
}

type itemType struct {
	URL   string      `yaml:"url"`
	Group []groupType `yaml:"group"`
}

// Template holds resources in yaml.
type Template struct {
	Items []itemType `yaml:"items"`
}

// ParseBytes parses bytes to Template.
func ParseBytes(b []byte) (*Template, error) {
	var out Template
	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ParseFile parses file to Template.
func ParseFile(filename string) (*Template, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseBytes(b)
}
