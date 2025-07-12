package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Val struct {
	ID       string `yaml:"id"`
	JSONPath string `yaml:"jsonpath"`
}

type Item struct {
	URL string `yaml:"url"`
	Val []Val  `yaml:"val"`
}
type Template struct {
	Template []Item
}

func ParseTemplate(filename string) (*Template, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var out Template
	if err = yaml.Unmarshal(f, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
