package tpl

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Group struct {
	ID       string `yaml:"id"`
	JSONPath string `yaml:"jsonpath"`
}

type Item struct {
	URL   string  `yaml:"url"`
	Group []Group `yaml:"group"`
}
type Template struct {
	Template []Item
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
