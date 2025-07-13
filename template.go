package tpl

import (
	"os"

	"gopkg.in/yaml.v3"
)

type groupType struct {
	ID       string `yaml:"id"`
	JSONPath string `yaml:"jsonpath"`
}

type itemType struct {
	URL   string      `yaml:"url"`
	Group []groupType `yaml:"group"`
}
type Template struct {
	Template []itemType
}

func ParseBytes(b []byte) (*Template, error) {
	var out Template
	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func ParseFile(filename string) (*Template, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseBytes(b)
}
