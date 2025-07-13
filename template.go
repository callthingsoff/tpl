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

// Template describes resources in yaml.
type Template struct {
	Template []itemType
}

// ParseBytes parses bytes to Template.
func ParseBytes(b []byte) (*Template, error) {
	var out Template
	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ParseBytes parses yaml file to Template.
func ParseFile(filename string) (*Template, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseBytes(b)
}
